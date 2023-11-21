package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/EphemSpirit/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T) Entry {
	acct := CreateRandomAccount(t)
	args := CreateEntryParams{
		AccountID: acct.ID,
		Amount: util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), args)
	
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, acct.ID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t);
}

func TestGetEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateRandomEntry(t)
	}

	args := ListEntriesParams{
		Limit: 2,
		Offset: 2,
	}

	entries, err := testQueries.ListEntries(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, entries, 2)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestUpdateEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)

	args := UpdateEntryParams{
		ID: entry1.ID,
		Amount: util.RandomMoney(),
	}

	entry2, err := testQueries.UpdateEntry(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, args.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry1.ID)

	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}