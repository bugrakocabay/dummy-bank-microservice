package db

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	TransactionID     string         `json:"transaction_id"`
	FromAccountID     string         `json:"from_account_id"`
	ToAccountID       string         `json:"to_account_id"`
	TransactionAmount int32          `json:"transaction_amount"`
	Description       sql.NullString `json:"description"`
}

// Account contains the response parameters of the account service
type Account struct {
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Balance   int32     `json:"balance"`
	Type      string    `json:"type"`
	AccountID string    `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transaction Transaction `json:"transaction"`
	FromAccount Account     `json:"from_account"`
	ToAccount   Account     `json:"to_account"`
}

// TransferTx Performs a money transfer from one account to the other.
// Creates a transfer record, adds account entries and updates accounts' balances within a single db transaction.
func (store *SQLStore) TransferTx(w http.ResponseWriter, ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			TransactionID:     store.createUUID(),
			FromAccountID:     arg.FromAccountID,
			ToAccountID:       arg.ToAccountID,
			Description:       arg.Description,
			TransactionAmount: arg.TransactionAmount,
		})
		if err != nil {
			log.Println(err)
			return err
		}

		result.FromAccount, err = store.updateAccountRequest(w, UpdatePayload{
			ID:      arg.FromAccountID,
			Balance: -arg.TransactionAmount,
		})
		if err != nil {
			log.Println(err)
			return err
		}

		result.ToAccount, err = store.updateAccountRequest(w, UpdatePayload{
			ID:      arg.ToAccountID,
			Balance: arg.TransactionAmount,
		})
		if err != nil {
			log.Println(err)
			return err
		}

		return nil
	})

	return result, err
}

func (store *SQLStore) createUUID() string {
	// Generate a new UUID
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		panic(err)
	}

	// Set the UUID version and variant bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0xbf) | 0x80

	// Convert the UUID to a string format
	return fmt.Sprintf("%x-%x-%x-%x-%x\n", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

type UpdatePayload struct {
	ID      string `json:"id" binding:"required,min=1"`
	Balance int32  `json:"balance" binding:"required"`
}

func (store *SQLStore) updateAccountRequest(w http.ResponseWriter, payload UpdatePayload) (Account, error) {
	jsonData, _ := json.Marshal(payload)

	request, err := http.NewRequest(http.MethodPut, "http://account-service/accounts/update", bytes.NewBuffer(jsonData))
	if err != nil {
		return Account{}, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return Account{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Account{}, err
	}

	maxBytes := 10485376 // 1mgb
	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody Account
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		return Account{}, err
	}

	return jsonResponseBody, nil
}
