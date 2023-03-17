package main

import (
	"github.com/dummy-bank-scripts/requests"
)

func main() {
	id := requests.CreateUser()
	accessToken := requests.Login(id)
	accountIDs := requests.GetAllAccounts(accessToken)
	randomAccountID := accountIDs[requests.RandomInt(0, int64(len(accountIDs)))]
	myAccountID := requests.CreateAccount(accessToken)

}
