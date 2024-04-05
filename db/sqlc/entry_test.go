package db

import (
	"context"
	"testing"

	"github.com/HeavenAQ/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	// create a dummy account
	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomBalance(),
	}

	// ensure nothing went wrong
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.UpdatedAt)
	return entry
}

func TestCreateEntry(t *testing.T) {
	entry := createRandomEntry(t)
	testQueries.DeleteEntry(context.Background(), entry.ID)
}

func TestGetEntry(t *testing.T) {
	// create a dummy entry
	entry1 := createRandomEntry(t)

	// get the entry and ensure it is the same as the one created
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	require.Equal(t, entry1.UpdatedAt, entry2.UpdatedAt)
	testQueries.DeleteEntry(context.Background(), entry1.ID)
}

func TestUpdateEntry(t *testing.T) {
	// create a dummy entry
	entry1 := createRandomEntry(t)
	arg := UpdateEntryParams{
		ID:        entry1.ID,
		AccountID: entry1.AccountID,
		Amount:    utils.RandomBalance(),
	}
	entry2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, arg.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	require.NotEqual(t, entry1.UpdatedAt, entry2.UpdatedAt)
	testQueries.DeleteEntry(context.Background(), entry1.ID)
}

func TestDeleteEntry(t *testing.T) {
	// create a dummy entry
	entry1 := createRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	// ensure the entry is deleted
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.Empty(t, entry2)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}
	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
