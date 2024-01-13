package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/EphemSpirit/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	user := CreateRandomUser(t)
	arg := CreateAccountParams{
		Owner: user.Username,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acctOne := CreateRandomAccount(t)	
	acctTwo, err := testQueries.GetAccount(context.Background(), acctOne.ID)
	require.NoError(t, err)
	require.NotEmpty(t, acctTwo)

	require.Equal(t, acctOne.Balance, acctTwo.Balance)
	require.Equal(t, acctOne.Owner, acctTwo.Owner)
	require.Equal(t, acctOne.ID, acctTwo.ID)
	require.Equal(t, acctOne.Currency, acctTwo.Currency)
	require.WithinDuration(t, acctOne.CreatedAt, acctTwo.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acct1 := CreateRandomAccount(t)

	arg := UpdateAccountParams{
		ID: acct1.ID,
		Balance: util.RandomMoney(),
	}

	acct2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, acct2)

	require.Equal(t, arg.Balance, acct2.Balance)
	require.Equal(t, acct1.Owner, acct2.Owner)
	require.Equal(t, acct1.ID, acct2.ID)
	require.Equal(t, acct1.Currency, acct2.Currency)
	require.WithinDuration(t, acct2.CreatedAt, acct2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acct1 := CreateRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acct1.ID)
	require.NoError(t, err)

	acct2, err := testQueries.GetAccount(context.Background(), acct1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acct2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = CreateRandomAccount(t)
	}

	args := ListAccountsParams{
		Owner: lastAccount.Owner,
		Limit: 5,
		Offset: 0,
	}

	accts, err := testQueries.ListAccounts(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, accts)

	for _, acct := range accts {
		require.NotEmpty(t, acct)
		require.Equal(t, lastAccount.Owner, acct.Owner)
	}
}