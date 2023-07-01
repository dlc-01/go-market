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
	m.Mutex.Lock()
	m.Mutex.Unlock()
	return mockS, nil

}

func (m TestStore) createUser(ctx context.Context, u *model.UserInfo) error {
	m.Mutex.Lock()
	ret := m.Called(ctx, u)
	var r0 error
	if ret != nil {
		if ret.Get(0) != nil {
			new := ret.Get(0).(*model.UserInfo)
			if u.Login == new.Login {
				return apperrors.NewConflict("username with login: " + u.Login)
			}
			m.Mutex.Unlock()
			return nil

		}
	}
	m.Mutex.Unlock()
	return r0
}

func (m TestStore) findByLogin(ctx context.Context, login *string) (*model.UserInfo, error) {
	m.Mutex.Lock()
	ret := m.Called(ctx, login)

	var r0 *model.UserInfo
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.UserInfo)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}
	m.Mutex.Unlock()
	return r0, r1
}
func (m TestStore) getAllOrdersByLogin(ctx context.Context, login *string) (*model.User, error) {
	m.Mutex.Lock()
	ret := m.Called(ctx, login)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}
	if len(r0.Orders) == 0 {
		return r0, apperrors.NewNoContent()
	}
	m.Mutex.Unlock()
	return r0, nil
}

func (m TestStore) addNewOrder(ctx context.Context, u *model.User) error {\
	m.Mutex.Lock()
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
	m.Mutex.Unlock()
	return r1
}

func (m TestStore) getBalanceWithdraw(ctx context.Context, login *string) (*model.BalanceResp, error) {
	m.Mutex.Lock()
	ret := m.Called(ctx, login)

	var r0 *model.BalanceResp
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.BalanceResp)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}
	m.Mutex.Unlock()
	return r0, r1
}

func (m TestStore) getUBalance(ctx context.Context, login *string) (float64, error) {
	m.Mutex.Lock()
	ret := m.Called(ctx, login)

	var r0 float64
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(float64)
	}
	m.Mutex.Unlock()
	return r0, nil
}

func (m TestStore) getAllWithdrawsByLogin(ctx context.Context, login *string) (*model.User, error) {
	m.Mutex.Lock()
	ret := m.Called(ctx, login)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}
	if len(r0.Withdraws) == 0 {
		m.Mutex.Unlock()
		return r0, apperrors.NewNoContent()
	}
	m.Mutex.Unlock()
	return r0, nil

}

func (m TestStore) addNewOderWithdraw(ctx context.Context, u *model.User) error {
	m.Mutex.Lock()
	m.Mutex.Unlock()
	return apperrors.NewAccepted()
}

var mockS UserStorage = TestStore{}
