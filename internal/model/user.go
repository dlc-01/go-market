package model

type User struct {
	Info      UserInfo
	Orders    []Order
	Withdraws []Withdraw
	Balance   float64
}

type UserInfo struct {
	Password string
	Login    string
}

type AuthReq struct {
	Username string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
