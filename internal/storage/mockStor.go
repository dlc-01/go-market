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
	ret := m.Called(ctx, login)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}
	if len(r0.Orders) == 0 {
		return r0, apperrors.NewNoContent()
	}

	return r0, nil
}

func (m TestStore) addNewOrder(ctx context.Context, u *model.User) error {
	ret := m.Called(ctx, u)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}
	var r1 error
	if r0.Info.Login == "" && len(r0.Orders) == 0 {
		r1 = apperrors.NewAccepted()
	} else if r0.Info.Login == u.Info.Login && r0.Orders[0].ID == u.Orders[0].ID {
		r1 = apperrors.NewStatusOK()
	} else if r0.Info.Login != u.Info.Login && r0.Orders[0].ID == u.Orders[0].ID {
		r1 = apperrors.NewConflict("")
	}

	return r1
}

func (m TestStore) getBalanceWithdraw(ctx context.Context, login *string) (*model.BalanceResp, error) {
	ret := m.Called(ctx, login)

	var r0 *model.BalanceResp
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.BalanceResp)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m TestStore) getUBalance(ctx context.Context, login *string) (float64, error) {
	ret := m.Called(ctx, login)

	var r0 float64
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(float64)
	}

	return r0, nil
}

func (m TestStore) getAllWithdrawsByLogin(ctx context.Context, login *string) (*model.User, error) {
	ret := m.Called(ctx, login)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}
	if len(r0.Withdraws) == 0 {
		return r0, apperrors.NewNoContent()
	}

	return r0, nil

}

func (m TestStore) addNewOderWithdraw(ctx context.Context, u *model.User) error {
	return apperrors.NewAccepted()
}

var mockS UserStorage = TestStore{}
