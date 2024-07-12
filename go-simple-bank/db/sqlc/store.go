package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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

// TransferTxParams is the params of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txtKey = struct{}{}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts balance within a single database transaction
/**
Sample SQL Query:
BEGIN;
	INSERT INTO transfers (from_account_id, to_account_id, amount) VALUES (1, 2, 10) RETURNING *;

	INSERT INTO entries (account_id, amount) VALUES (1, -10) RETURNING *;
	INSERT INTO entries (account_id, amount) VALUES (2, 10) RETURNING *;

	SELECT * FROM accounts WHERE id = 1 FOR UPDATE;
	UPDATE accounts SET balance = 90 WHERE id = 1 RETURNING *;

	SELECT * FROM accounts WHERE id = 2 FOR UPDATE;
	UPDATE accounts SET balance = 110 WHERE id = 2 RETURNING *;
ROLLBACK;
*/
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txtKey)

		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// get accounts -> update account's balance
		//fmt.Println(txName, "get from account")
		//account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		//if err != nil {
		//	return err
		//}
		//
		//fmt.Println(txName, "update from account balance")
		//result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		//	ID:      arg.FromAccountID,
		//	Balance: account1.Balance - arg.Amount,
		//})
		//if err != nil {
		//	return err
		//}
		//
		//fmt.Println(txName, "get to account")
		//account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		//if err != nil {
		//	return err
		//}
		//
		//fmt.Println(txName, "update to account balance")
		//result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		//	ID:      arg.ToAccountID,
		//	Balance: account2.Balance + arg.Amount,
		//})
		//if err != nil {
		//	return err
		//}

		// get accounts -> update account's balance
		// This one should arrange the order of the UPDATE statement in transaction
		// to prevent deadlock
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return
}
