package storage

import (
	"context"
	"fmt"
	"github.com/dlc/go-market/internal/config"
	"github.com/dlc/go-market/internal/model"
)

type UserStorage interface {
	CreateStor(ctx context.Context, cfg *config.ServerConfig) (UserStorage, error)
	CreateUser(ctx context.Context, u *model.UserInfo) error
	FindByLogin(ctx context.Context, login *string) (*model.UserInfo, error)
	AddNewOrder(ctx context.Context, u *model.User) error
	GetAllOrdersByLogin(ctx context.Context, login *string) (*model.User, error)
	GetBalanceWithdraw(ctx context.Context, login *string) (*model.BalanceResp, error)
	AddNewOderWithdraw(ctx context.Context, u *model.User) error
	GetAllWithdrawsByLogin(ctx context.Context, login *string) (*model.User, error)
	CollectOrders(ctx context.Context) ([]model.Order, error)
	UpdateOrders(ctx context.Context, order model.Order) error
	Close(ctx context.Context)
}

var stor UserStorage

func InitDBStorage(ctx context.Context, cfg *config.ServerConfig) error {
	var err error
	stor, err = db.CreateStor(ctx, cfg)
	if err != nil {
		return fmt.Errorf("cannot init storage: %w", err)
	}
	return nil
}

func InitTestStorage(m *TestStore) {
	stor = m
}

func CreateUser(ctx context.Context, u *model.UserInfo) error {
	return stor.CreateUser(ctx, u)
}

func FindByLogin(ctx context.Context, login *string) (*model.UserInfo, error) {
	return stor.FindByLogin(ctx, login)
}

func AddNewOrder(ctx context.Context, u *model.User) error {
	return stor.AddNewOrder(ctx, u)
}

func GetAllOrdersByLogin(ctx context.Context, login *string) (*model.User, error) {
	return stor.GetAllOrdersByLogin(ctx, login)
}

func GetUserBalanceWithdraw(ctx context.Context, login *string) (*model.BalanceResp, error) {
	return stor.GetBalanceWithdraw(ctx, login)
}

func AddNewOderWithdraw(ctx context.Context, u *model.User) error {
	return stor.AddNewOderWithdraw(ctx, u)
}

func GetAllWithdrawsByLogin(ctx context.Context, login *string) (*model.User, error) {
	return stor.GetAllWithdrawsByLogin(ctx, login)
}

func CollectOrders(ctx context.Context) ([]model.Order, error) {
	return stor.CollectOrders(ctx)
}

func UpdateOrders(ctx context.Context, order model.Order) error {

	return stor.UpdateOrders(ctx, order)
}

func Close(ctx context.Context) {
	stor.Close(ctx)
}
