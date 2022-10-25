package db

import (
	"context"
	"log"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

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

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

/* TransferTx perform transfer from one account to another account
It creates a transfer recored, entries and update the account balance
with single database transacton.
*/
func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := crdbpgx.ExecuteTx(ctx, store.db, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) error {
		var err error

		result.Transfer, err = store.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = store.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = store.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}
		// Update Account Balance
		if args.FromAccountID < args.ToAccountID {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, store.Queries, args.FromAccountID, -args.Amount, args.ToAccountID, args.Amount)
			if err != nil {
				log.Printf("AddMoney error %v", err)
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = AddMoney(ctx, store.Queries, args.ToAccountID, args.Amount, args.FromAccountID, -args.Amount)
			if err != nil {
				log.Printf("AddMoney error %v", err)
				return err
			}

		}
		// Print transaction result
		log.Printf("TransferTx result %v", result)

		return nil
	})

	if err != nil {
		log.Println("Error while creating transaction", err)
	}

	return result, err

}

func AddMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (acount1 Account, acount2 Account, err error) {
	acount1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	acount2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return
}
