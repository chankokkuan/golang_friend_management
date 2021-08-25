package port

import (
	"context"
	"friend-management/internal/core/domain"
)

type UserRepository interface {
	ConnectionCheck(context.Context) bool
	CreateUser(context.Context, domain.CreateUserRequest) (*domain.User, error)
	GetUsers(context.Context, domain.UserQuery) (*domain.MetaUsers, error)
	GetUser(context.Context, string) (*domain.User, error)
	UpdateUser(context.Context, domain.UpdateUserRequest) (*domain.User, error)
	AddFriend(context.Context, domain.AddFriendRequest) error
}
