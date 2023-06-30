package storage

import (
	"context"
	"github.com/dlc/go-market/internal/config"
	"github.com/dlc/go-market/internal/model"
	"github.com/dlc/go-market/internal/model/apperrors"

	"github.com/stretchr/testify/mock"
)

type MockStor struct {
	mock.Mock
}

func (m MockStor) createStor(ctx context.Context, cfg *config.ServerConfig) (UserStorage, error) {
	return mockS, nil
}

func (m MockStor) createUser(ctx context.Context, u *model.UserInfo) error {
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

func (m MockStor) findByLogin(ctx context.Context, login *string) (*model.UserInfo, error) {
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
func (m MockStor) getAllOrdersByLogin(ctx context.Context, login *string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockStor) addNewOrder(ctx context.Context, u *model.User) error {
	return nil
}

func (m MockStor) getBalanceWithdraw(ctx context.Context, login *string) (*model.BalanceResp, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockStor) getUBalance(ctx context.Context, login *string) (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockStor) getAllWithdrawsByLogin(ctx context.Context, login *string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockStor) addNewOderWithdraw(ctx context.Context, u *model.User) error {
	//TODO implement me
	panic("implement me")
}

var mockS UserStorage = MockStor{}
