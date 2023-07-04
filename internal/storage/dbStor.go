package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/dlc/go-market/internal/config"
	"github.com/dlc/go-market/internal/hash"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/model/apperrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbStor struct {
	*pgxpool.Pool
}

const migrationUserTable = "CREATE TABLE IF NOT EXISTS market_users ( username VARCHAR NOT NULL PRIMARY KEY,  password VARCHAR NOT NULL);"
const migrationOrdersTable = "CREATE TABLE IF NOT EXISTS market_orders (id VARCHAR NOT NULL PRIMARY KEY,client VARCHAR NOT NULL, status VARCHAR, accrual NUMERIC(10,2) NOT NULL, time_created timestamp NOT NULL, FOREIGN KEY (client) REFERENCES market_users(username) ON DELETE CASCADE );"
const migrationBalanceTable = "CREATE TABLE IF NOT EXISTS market_ubalance (client VARCHAR NOT NULL UNIQUE, balance NUMERIC(10,2) NOT NULL, FOREIGN KEY (client) REFERENCES market_users(username)  ON DELETE CASCADE); "
const migrationWithdrawsTable = "CREATE TABLE IF NOT EXISTS market_withdraws (client VARCHAR NOT NULL , order_id VARCHAR NOT NULL UNIQUE , amount NUMERIC(10,2) NOT NULL,time_created TIMESTAMP NOT NULL, FOREIGN KEY (client) REFERENCES market_users(username) ON DELETE CASCADE, FOREIGN KEY (order_id) REFERENCES market_orders(id) ON DELETE CASCADE, PRIMARY KEY (client, order_id));"

var db UserStorage = &dbStor{}

func (d dbStor) CreateStor(ctx context.Context, cfg *config.ServerConfig) (UserStorage, error) {
	var err error

	d.Pool, err = pgxpool.New(ctx, cfg.DBAddress)
	if err != nil {
		return &d, fmt.Errorf("cannot connect to db: %w ", err)
	}
	if err = d.Ping(ctx); err != nil {
		return &d, fmt.Errorf("error while try to ping db: %w", err)
	}
	if _, err = d.migration(ctx); err != nil {
		return nil, fmt.Errorf("cannot create migration %w", err)
	}
	return &d, nil

}
func (d dbStor) migration(ctx context.Context) (UserStorage, error) {
	if _, err := d.Exec(ctx, migrationUserTable); err != nil {
		return nil, fmt.Errorf("cannot create init migration user table %w", err)
	}

	if _, err := d.Exec(ctx, migrationOrdersTable); err != nil {
		return nil, fmt.Errorf("cannot create migration orders table %w", err)
	}

	if _, err := d.Exec(ctx, migrationBalanceTable); err != nil {
		return nil, fmt.Errorf("cannot create  migration balance table %w", err)
	}

	if _, err := d.Exec(ctx, migrationWithdrawsTable); err != nil {
		return nil, fmt.Errorf("cannot create migration withdraws table %w", err)
	}

	return &d, nil

}

func (d dbStor) CreateUser(ctx context.Context, u *model.UserInfo) error {
	var err error
	u.Password, err = hash.HashPassword(u.Password)
	if err != nil {
		return err
	}
	err = d.QueryRow(ctx, "INSERT INTO market_users (username, password) VALUES ($1, $2)", u.Login, u.Password).Scan()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "market_users_pkey" {
				return apperrors.NewConflict("username with " + u.Login)
			}
			return err
		}
	}

	if err := d.createBalance(ctx, u.Login); err != nil {
		return fmt.Errorf("error while creating balance: %w", err)
	}

	return nil
}

func (d dbStor) createBalance(ctx context.Context, login string) error {
	err := d.QueryRow(ctx, "INSERT INTO market_ubalance (client, balance) VALUES ($1, $2)", login, 0).Scan()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return err
		}

	}
	return nil
}

func (d dbStor) FindByLogin(ctx context.Context, login *string) (*model.UserInfo, error) {
	resU := model.UserInfo{}
	err := d.QueryRow(ctx, "SELECT * FROM market_users WHERE username=$1", login).Scan(&resU.Login, &resU.Password)
	if err != nil {
		if pgx.ErrNoRows == err {
			return &resU, apperrors.NewUnauthorized("not found username")
		}
		return &resU, fmt.Errorf("error while scan response from db : %w", err)
	}
	return &resU, nil
}
func (d dbStor) AddNewOrder(ctx context.Context, u *model.User) error {
	if err := d.CheckUniqOrder(ctx, u); err != nil {
		return err
	}
	order := u.Orders[0]
	err := d.QueryRow(ctx, "INSERT INTO market_orders (id, client, status, accrual, time_created) VALUES ($1, $2, $3, $4, $5)", order.ID, u.Info.Login, order.Status, order.Accrual, order.TimeCreated).Scan()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("eroror while creating new order: %w", err)
		}
	}

	return apperrors.NewAccepted()

}
func (d dbStor) CheckUniqOrder(ctx context.Context, u *model.User) error {
	order := u.Orders[0]
	var clientFromDB string
	err := d.QueryRow(ctx, "SELECT client FROM market_orders WHERE id = $1", order.ID).Scan(&clientFromDB)
	if err != nil {
		if pgx.ErrNoRows != err {
			return fmt.Errorf("error while checking orders : %w", err)
		} else {
			return nil
		}
	}

	if clientFromDB != u.Info.Login {
		return apperrors.NewConflict("order with other login ")
	} else {
		return apperrors.NewStatusOK()
	}
}

func (d dbStor) GetAllOrdersByLogin(ctx context.Context, login *string) (*model.User, error) {
	u := model.User{Info: model.UserInfo{Login: *login}}
	orders := make([]model.Order, 0)

	col, err := d.Query(ctx, "SELECT id, status, accrual, time_created FROM market_orders WHERE client = $1 ORDER BY time_created ", u.Info.Login)
	if err != nil {
		return &u, fmt.Errorf("error while sending query to db : %w", err)
	}

	for col.Next() {
		order := model.Order{}
		if err := col.Scan(&order.ID, &order.Status, &order.Accrual, &order.TimeCreated); err != nil {
			return &u, fmt.Errorf("error file scanning resp from db: %w", err)
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return &u, apperrors.NewNoContent()
	}

	u.Orders = orders
	return &u, nil
}

func (d dbStor) GetBalanceWithdraw(ctx context.Context, login *string) (*model.BalanceResp, error) {
	b := &model.BalanceResp{}
	var err error
	if b.Balance, err = d.GetUBalance(ctx, login); err != nil {
		return b, fmt.Errorf("error while get balance: %w", err)
	}
	if b.Sum, err = d.GetSumOfWithdraws(ctx, login); err != nil {
		return b, fmt.Errorf("error while get sum withdraws: %w", err)
	}
	return b, err

}

func (d dbStor) GetUBalance(ctx context.Context, login *string) (float64, error) {
	var balance float64

	err := d.QueryRow(ctx, "SELECT balance FROM  market_ubalance WHERE client = $1", *login).Scan(&balance)
	if err != nil {
		return balance, fmt.Errorf("error while scanning query to db : %w", err)
	}

	return balance, nil
}

func (d dbStor) GetSumOfWithdraws(ctx context.Context, login *string) (float64, error) {
	var withdraw float64

	//TODO обработать ошибку если сканим NULL
	err := d.QueryRow(ctx, "SELECT SUM(amount) FROM  market_withdraws WHERE client = $1", *login).Scan(&withdraw)
	if err != nil {
		return 0, nil
	}

	return withdraw, nil
}

func (d dbStor) AddNewOderWithdraw(ctx context.Context, u *model.User) error {
	if err := d.AddNewOrder(ctx, u); err != nil {
		if apperrors.Status(err) != 202 {
			return err
		}
	}
	withdraw := u.Withdraws[0]
	err := d.QueryRow(ctx, "INSERT INTO market_withdraws (client, order_id, amount, time_created) VALUES ($1, $2, $3, $4)", u.Info.Login, withdraw.Order, withdraw.Sum, withdraw.TimeCreated).Scan()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("eroror while creating new order: %w", err)
		}
	}
	if err := d.UpdateUBalance(ctx, u); err != nil {
		return fmt.Errorf("error while updating balance")
	}

	return apperrors.NewAccepted()

}
func (d dbStor) UpdateUBalance(ctx context.Context, u *model.User) error {

	err := d.QueryRow(ctx, "UPDATE market_ubalance SET balance = $1 WHERE client = $2", u.Balance, u.Info.Login).Scan()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("eroror while updating balance: %w", err)
		}
	}

	return nil
}

func (d dbStor) GetAllWithdrawsByLogin(ctx context.Context, login *string) (*model.User, error) {
	u := model.User{Info: model.UserInfo{Login: *login}}
	withdraws := make([]model.Withdraw, 0)

	col, err := d.Query(ctx, "SELECT order_id, amount, time_created FROM market_withdraws WHERE client = $1 ORDER BY time_created ASC", u.Info.Login)
	if err != nil {
		return &u, fmt.Errorf("error while sending query to db : %w", err)
	}

	for col.Next() {
		withdraw := model.Withdraw{}
		if err := col.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.TimeCreated); err != nil {
			return &u, fmt.Errorf("error file scanning resp from db: %w", err)
		}
		withdraws = append(withdraws, withdraw)
	}

	if len(withdraws) == 0 {
		return &u, apperrors.NewNoContent()
	}

	u.Withdraws = withdraws
	return &u, nil
}
func (d dbStor) CollectOrders(ctx context.Context) ([]model.Order, error) {

	orders := make([]model.Order, 0)

	col, err := d.Query(ctx, "SELECT id, status FROM market_orders WHERE status = $1 or status = $2 ", model.NEW, model.PROCESSED)
	if err != nil {
		return nil, fmt.Errorf("error while sending query to db : %w", err)
	}

	for col.Next() {
		order := model.Order{}
		if err := col.Scan(&order.ID, &order.Status); err != nil {
			return orders, fmt.Errorf("error file scanning resp from db: %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil

}

func (d dbStor) UpdateOrders(ctx context.Context, order model.Order) error {

	err := d.QueryRow(ctx, "UPDATE market_orders SET status = $1, accrual = $2 WHERE id = $3", order.Status, order.Accrual, order.ID).Scan()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("error while sending query to db: %w", err)
		}
	}
	return d.updateBalanceWithAccrual(ctx, order.Accrual, order.ID)
}

func (d dbStor) updateBalanceWithAccrual(ctx context.Context, accrual float64, order string) error {

	err := d.QueryRow(ctx, "UPDATE market_ubalance SET balance = balance + $1 WHERE client = (SELECT client FROM market_orders WHERE id = $2)", accrual, order).Scan()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("eroror while updating balance with accrual: %w", err)
		}
	}

	return nil
}

func (db dbStor) Close(ctx context.Context) {
	db.Close(ctx)
}
