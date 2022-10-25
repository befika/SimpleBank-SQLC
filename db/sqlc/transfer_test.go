package db

import (
	"context"
	"simple_bank/util"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T) Transfer {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	args := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMony(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	t1 := CreateRandomTransfer(t)
	t2, err := testQueries.GetTransfer(context.Background(), t1.ID)
	require.NoError(t, err)
	require.Equal(t, t1.ID, t2.ID)
	require.Equal(t, t1.FromAccountID, t2.FromAccountID)
	require.Equal(t, t1.ToAccountID, t2.ToAccountID)
	require.Equal(t, t1.CreatedAt, t2.CreatedAt)
}

func TestUpdateTransfer(t *testing.T) {
	t1 := CreateRandomTransfer(t)
	args := UpdateTransferParams{
		ID:     t1.ID,
		Amount: t1.Amount,
	}
	t2, err := testQueries.UpdateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, t2)
	require.Equal(t, t1.ID, t2.ID)
	require.Equal(t, t1.Amount, t2.Amount)
	require.Equal(t, t1.CreatedAt, t2.CreatedAt)
}

func TestDeleteTransfer(t *testing.T) {
	t1 := CreateRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), t1.ID)
	require.NoError(t, err)

	t2, err := testQueries.GetTransfer(context.Background(), t1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, t2)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t)
	}
	param := ListTransferParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQueries.ListTransfer(context.Background(), param)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, t1 := range transfers {
		require.NotEmpty(t, t1)
	}
}
