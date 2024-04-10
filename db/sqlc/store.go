package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db DBTX) *Store {
	return &Store{
		db:      db.(*pgxpool.Pool),
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}
	return tx.Commit(ctx)
}

// Info for transfer transactions
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// create transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// create an entry of from account
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// create an entry of to account
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// transfer money between accounts
		result.FromAccount, result.ToAccount, err = transferMoney(ctx, q, arg)
		if err != nil {
			return err
		}
		return nil
	})

	return result, err
}

func transferMoney(
	ctx context.Context,
	q *Queries,
	arg TransferTxParams,
) (fromAccount Account, toAccount Account, err error) {
	fromAccountID, toAccountID := arg.FromAccountID, arg.ToAccountID
	amount := arg.Amount
	if fromAccountID > toAccountID {
		fromAccountID, toAccountID = toAccountID, fromAccountID
		amount = -amount
	}

	// update sender's account balance
	fromAccount, err = q.AdjustAccountBalance(ctx, AdjustAccountBalanceParams{
		ID:     fromAccountID,
		Amount: -amount,
	})
	if err != nil {
		return
	}

	// update receiver's account balance
	toAccount, err = q.AdjustAccountBalance(ctx, AdjustAccountBalanceParams{
		ID:     toAccountID,
		Amount: amount,
	})
	if err != nil {
		return
	}
	return
}
