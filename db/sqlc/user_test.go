package db

import (
	"context"
	"testing"
	"time"

	"github.com/EphemSpirit/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username: util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	userOne := CreateRandomUser(t)	
	userTwo, err := testQueries.GetUser(context.Background(), userOne.Username)
	require.NoError(t, err)
	require.NotEmpty(t, userTwo)

	require.Equal(t, userOne.Username, userTwo.Username)
	require.Equal(t, userOne.HashedPassword, userTwo.HashedPassword)
	require.Equal(t, userOne.Email, userTwo.Email)
	require.Equal(t, userOne.FullName, userTwo.FullName)
	require.WithinDuration(t, userOne.CreatedAt, userTwo.CreatedAt, time.Second)
}