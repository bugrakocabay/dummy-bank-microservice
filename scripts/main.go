package main

import (
	"fmt"
	"github.com/dummy-bank-scripts/requests"
)

func main() {
	// This is to create dummy data, that will fill daily-report service
	/*for i := 0; i < 10; i++ {
		id := requests.CreateUser()
		accessToken := requests.Login(id)
		myAccountID := requests.CreateAccount(accessToken)
		err := requests.AddBalance(accessToken, myAccountID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created user: %d\n", i+1)
	}*/
	for i := 0; i < 100; i++ {
		email := requests.CreateUser()
		accessToken := requests.Login(email)
		accountIDs := requests.GetAllAccounts(accessToken)
		myAccountID := requests.CreateAccount(accessToken)
		randomAccountID := accountIDs[requests.RandomInt(0, int64(len(accountIDs)-1))]
		err := requests.AddBalance(accessToken, myAccountID)
		if err != nil {
			panic(err)
		}

		err = requests.CreateTransaction(accessToken, myAccountID, randomAccountID)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%d done!\n", i+1)
	}
}
