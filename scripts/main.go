package main

import (
	"fmt"
	"github.com/dummy-bank-scripts/requests"
)

func main() {
	for i := 0; i < 100; i++ {
		id := requests.CreateUser()
		accessToken := requests.Login(id)
		accountIDs := requests.GetAllAccounts(accessToken)
		randomAccountID := accountIDs[requests.RandomInt(0, int64(len(accountIDs)))]
		myAccountID := requests.CreateAccount(accessToken)
		err := requests.AddBalance(accessToken, myAccountID)
		if err != nil {
			panic(err)
		}

		err = requests.CreateTransaction(accessToken, myAccountID, randomAccountID)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%d done!", i+1)
	}
}
