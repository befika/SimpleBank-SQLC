package db

import (
	"context"
	"simple_bank/util"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T) Entry {
	account := CreateRandomAccount(t)

	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMony(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
}

func TestUpdateEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	args := UpdateEntryParams{
		ID:     entry1.ID,
		Amount: entry1.Amount,
	}
	entry2, err := testQueries.UpdateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
}

func TestDeleteEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	err := testQueries.DeleteAccount(context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetAccount(context.Background(), entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomEntry(t)
	}
	param := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListEntries(context.Background(), param)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, account := range entries {
		require.NotEmpty(t, account)
	}
}
