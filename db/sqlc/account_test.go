package db

import (
	"context"
	"testing"

	"github.com/HeavenAQ/simple-bank/utils"
	"github.com/stretchr/testify/require"
)

// Private Functions
func createRandomAccount(t *testing.T) Account {
	// create a dummy account
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomBalance(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	// ensure nothing went wrong
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	// ensure auto-generated fields are not empty
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	account := createRandomAccount(t)
	testQueries.DeleteAccount(context.Background(), account.ID)
}

func TestGetAccount(t *testing.T) {
	// create a dummy account
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	// ensure the account fetched is the same as the one created
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
	require.Equal(t, account1.UpdatedAt, account2.UpdatedAt)
	testQueries.DeleteAccount(context.Background(), account1.ID)
}

func TestUpdateAccount(t *testing.T) {
	// create a dummy account
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomBalance(),
	}

	// update the account and ensure nothing went wrong
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	// ensure the account fetched is the same as the one created
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
	require.NotEqual(t, account1.UpdatedAt, account2.UpdatedAt)
	testQueries.DeleteAccount(context.Background(), account1.ID)
}

func TestDeleteAccount(t *testing.T) {
	// create a dummy account
	account1 := createRandomAccount(t)

	// delete the account and ensure nothing went wrong
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	// ensure the account is deleted
	account, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.Empty(t, account)
}

func TestListAccounts(t *testing.T) {
	accountsCreated := []int64{}
	for i := 0; i < 10; i++ {
		accountsCreated = append(accountsCreated, createRandomAccount(t).ID)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	// fetch 5 accounts
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	// ensure all accounts are not empty
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

	// clean up
	for _, accountID := range accountsCreated {
		testQueries.DeleteAccount(context.Background(), accountID)
	}
}
