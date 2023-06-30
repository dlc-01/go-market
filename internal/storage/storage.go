package storage

import (
	"context"
	"fmt"
	"github.com/dlc/go-market/internal/config"
	"github.com/dlc/go-market/internal/model"
)

type UserStorage interface {
	createStor(ctx context.Context, cfg *config.ServerConfig) (UserStorage, error)
	createUser(ctx context.Context, u *model.UserInfo) error
	findByLogin(ctx context.Context, login *string) (*model.UserInfo, error)
	addNewOrder(ctx context.Context, u *model.User) error
	getAllOrdersByLogin(ctx context.Context, login *string) (*model.User, error)
	getUBalance(ctx context.Context, login *string) (float64, error)
	getBalanceWithdraw(ctx context.Context, login *string) (*model.BalanceResp, error)
	addNewOderWithdraw(ctx context.Context, u *model.User) error
	getAllWithdrawsByLogin(ctx context.Context, login *string) (*model.User, error)
}

var stor UserStorage

func InitDBStorage(ctx context.Context, cfg *config.ServerConfig) error {
	var err error
	stor, err = db.createStor(ctx, cfg)
	if err != nil {
		return fmt.Errorf("cannot init storage: %w", err)
	}
	return nil
}

func InitTestStorage(m *TestStore) {
	stor = m
}

func CreateUser(ctx context.Context, u *model.UserInfo) error {
	return stor.createUser(ctx, u)
}

func FindByLogin(ctx context.Context, login *string) (*model.UserInfo, error) {
	return stor.findByLogin(ctx, login)
}

func AddNewOrder(ctx context.Context, u *model.User) error {
	return stor.addNewOrder(ctx, u)
}

func GetAllOrdersByLogin(ctx context.Context, login *string) (*model.User, error) {
	return stor.getAllOrdersByLogin(ctx, login)
}

func GetUserBalanceWithdraw(ctx context.Context, login *string) (*model.BalanceResp, error) {
	return stor.getBalanceWithdraw(ctx, login)
}

func GetUserBalance(ctx context.Context, login *string) (float64, error) {
	return stor.getUBalance(ctx, login)
}

func AddNewOderWithdraw(ctx context.Context, u *model.User) error {
	return stor.addNewOderWithdraw(ctx, u)
}

func GetAllWithdrawsByLogin(ctx context.Context, login *string) (*model.User, error) {
	return stor.getAllWithdrawsByLogin(ctx, login)
}
