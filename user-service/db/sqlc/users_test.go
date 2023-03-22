package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		UserID:    RandomString(5),
		Email:     RandomEmail(),
		Password:  RandomString(10),
		Firstname: RandomString(5),
		Lastname:  RandomString(6),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.UserID, arg.UserID)
	require.Equal(t, user.Firstname, arg.Firstname)
	require.Equal(t, user.Lastname, arg.Lastname)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.Password, arg.Password)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Firstname, user2.Firstname)
	require.Equal(t, user1.Lastname, user2.Lastname)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestGetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Firstname, user2.Firstname)
	require.Equal(t, user1.Lastname, user2.Lastname)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserPassword(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUserPasswordParams{
		UserID:      user1.UserID,
		NewPassword: RandomString(10),
	}
	err := testQueries.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserID, user2.UserID)
	require.Equal(t, user1.Firstname, user2.Firstname)
	require.Equal(t, user1.Lastname, user2.Lastname)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, arg.NewPassword, user2.Password)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}