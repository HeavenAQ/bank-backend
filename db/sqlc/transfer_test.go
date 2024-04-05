package db

import (
	"context"
	"testing"

	"github.com/HeavenAQ/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	// create dummy accounts
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// create a dummy transfer
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        utils.RandomBalance(),
	}

	// create the transfer
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	// ensure values are correct
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	// ensure auto-generated fields are not empty
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	require.NotZero(t, transfer.UpdatedAt)
	return transfer
}

func TestCreatedTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)
	testQueries.DeleteTransfer(context.Background(), transfer.ID)
}

func TestGetTransfer(t *testing.T) {
	// create a dummy transfer
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	// ensure the transfer fetched is the same as the one created
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
	require.Equal(t, transfer1.UpdatedAt, transfer2.UpdatedAt)
	testQueries.DeleteTransfer(context.Background(), transfer1.ID)
}

func TestUpdateTransfer(t *testing.T) {
	// create a dummy transfer
	transfer1 := createRandomTransfer(t)
	arg := UpdateTransferParams{
		ID:            transfer1.ID,
		Amount:        utils.RandomBalance(),
		FromAccountID: transfer1.FromAccountID,
		ToAccountID:   transfer1.ToAccountID,
	}

	// update the transfer
	transfer2, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	// ensure the transfer is updated
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, arg.Amount, transfer2.Amount)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
	require.NotEqual(t, transfer1.UpdatedAt, transfer2.UpdatedAt)

	// delete the transfer
	testQueries.DeleteTransfer(context.Background(), transfer1.ID)
}

func TestDeleteTransfer(t *testing.T) {
	// create a dummy transfer
	transfer1 := createRandomTransfer(t)

	// delete the transfer
	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	// ensure the transfer is deleted
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.Empty(t, transfer2)
}

func TestListTransfer(t *testing.T) {
	// create 5 dummy transfers
	var transferIDs []int64
	for i := 0; i < 5; i++ {
		transfer := createRandomTransfer(t)
		transferIDs = append(transferIDs, transfer.ID)
	}

	// list the transfers
	arg := ListTransfersParams{
		Limit:  5,
		Offset: 0,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	// ensure the list is correct
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}

	// delete the transfers
	for _, id := range transferIDs {
		testQueries.DeleteTransfer(context.Background(), id)
	}
}
