package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// Run concurent transfer transactons
	n := 1
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// Check result
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check from entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// Check to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check Accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// Check Account Balance
		deff1 := account1.Balance - fromAccount.Balance
		deff2 := toAccount.Balance - account2.Balance
		require.Equal(t, deff1, deff2)
		require.True(t, deff1 > 0)
		require.True(t, deff1%amount == 0)

		k := int(deff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the final updated balance
	updateaccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateaccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updateaccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateaccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// Run concurent transfer transactons
	n := 10000
	amount := int64(10)

	errs := make(chan error)
	transactons := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			store := NewStore(testDb)
			transaction, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
			transactons <- transaction
		}()
	}

	// Check result
	for i := 0; i < n; i++ {
		err := <-errs
		transaction := <-transactons
		require.NoError(t, err)
		require.NotEmpty(t, transaction)

	}

	// Check the final updated balance
	updateaccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateaccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updateaccount1.Balance)
	require.Equal(t, account2.Balance, updateaccount2.Balance)

}
