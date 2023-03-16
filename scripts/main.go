package main

import (
	"fmt"

	"github.com/dummy-bank-scripts/handlers"
)

func main() {
	id := handlers.CreateUser()
	fmt.Println(id)
}
