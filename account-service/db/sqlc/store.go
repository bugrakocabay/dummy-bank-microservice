package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
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
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
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
	TransactionAmount float64        `json:"transaction_amount"`
	Description       sql.NullString `json:"description"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transaction Transaction `json:"transaction"`
	FromAccount Account     `json:"from_account"`
	ToAccount   Account     `json:"to_account"`
}

// TransferTx Performs a money transfer from one account to the other.
// Creates a transfer record, adds account entries and updates accounts' balances within a single db transaction.
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		commission := float64(arg.TransactionAmount) * 0.03
		result.Transaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			TransactionID:     store.createUUID(),
			FromAccountID:     arg.FromAccountID,
			ToAccountID:       arg.ToAccountID,
			Description:       arg.Description,
			TransactionAmount: arg.TransactionAmount,
			Commission:        commission,
		})
		if err != nil {
			log.Println(err)
			return err
		}
		moneyToBeTransferred := arg.TransactionAmount - commission

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, -arg.TransactionAmount, moneyToBeTransferred)
		} else {
			result.ToAccount, result.FromAccount, err = AddMoney(ctx, q, arg.ToAccountID, arg.FromAccountID, moneyToBeTransferred, -arg.TransactionAmount)
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
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func AddMoney(
	ctx context.Context,
	q *Queries,
	accountID1,
	accountID2 string,
	amount1,
	amount2 float64,
) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		AccountID: accountID1,
		Amount:    amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		AccountID: accountID2,
		Amount:    amount2,
	})
	return
}
