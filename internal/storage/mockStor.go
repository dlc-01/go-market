package storage

import (
	"context"
	"github.com/dlc/go-market/internal/config"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/model/apperrors"
	"sync"

	"github.com/stretchr/testify/mock"
)

type TestStore struct {
	mock.Mock
	sync.Mutex
}

func (m TestStore) createStor(ctx context.Context, cfg *config.ServerConfig) (UserStorage, error) {
	return mockS, nil
}

func (m TestStore) createUser(ctx context.Context, u *model.UserInfo) error {
	ret := m.Called(ctx, u)
	var r0 error
	if ret != nil {
		if ret.Get(0) != nil {
			new := ret.Get(0).(*model.UserInfo)
			if u.Login == new.Login {
				return apperrors.NewConflict("username with login: " + u.Login)
			}
			return nil

		}
	}

	return r0
}

func (m TestStore) findByLogin(ctx context.Context, login *string) (*model.UserInfo, error) {
	ret := m.Called(ctx, login)

	var r0 *model.UserInfo
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.UserInfo)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
func (m TestStore) getAllOrdersByLogin(ctx context.Context, login *string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m TestStore) addNewOrder(ctx context.Context, u *model.User) error {
	return nil
}

func (m TestStore) getBalanceWithdraw(ctx context.Context, login *string) (*model.BalanceResp, error) {
	//TODO implement me
	panic("implement me")
}

func (m TestStore) getUBalance(ctx context.Context, login *string) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (m TestStore) getAllWithdrawsByLogin(ctx context.Context, login *string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m TestStore) addNewOderWithdraw(ctx context.Context, u *model.User) error {
	//TODO implement me
	panic("implement me")
}

var mockS UserStorage = TestStore{}
