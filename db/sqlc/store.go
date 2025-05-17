package db

import (
	"context"
	"database/sql"
)

type Store interface {
	Querier
	TransferTx(context context.Context, arg CreateTransferParams) (result TransferTxResult, err error)
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

func (s *SQLStore) execTx(context context.Context, queryFunction func(*Queries) error) (err error) {

	tx, err := s.db.BeginTx(context, &sql.TxOptions{})
	if err != nil {
		return
	}

	if err = queryFunction(New(tx)); err != nil {
		if err = tx.Rollback(); err != nil {
			return
		}
		return
	}

	return tx.Commit()
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *SQLStore) TransferTx(context context.Context, arg CreateTransferParams) (result TransferTxResult, err error) {
	err = s.execTx(context, func(q *Queries) (err error) {
		result.Transfer, err = q.CreateTransfer(context, arg)
		if err != nil {
			return
		}

		result.FromEntry, err = q.CreateEntry(context, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return
		}
		result.ToEntry, err = q.CreateEntry(context, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addBalance(context, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addBalance(context, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return
	})
	return
}

func addBalance(context context.Context, q *Queries, accountId1 int64, amount1 int64, accountId2 int64, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(context, AddAccountBalanceParams{
		ID:     accountId1,
		Amount: amount1,
	},
	)
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(context, AddAccountBalanceParams{
		ID:     accountId2,
		Amount: amount2,
	},
	)
	return
}
