package services

import (
	"context"

	"github.com/kaa-dan/clean-architecture-go/internal/domain/entities"
)

type UserService interface {
	SignUp(ctx context.Context, req *entities.SignUpRequest) (*entities.AuthResponse, error)
	SignIn(ctx context.Context, req *entities.SignInRequest) (*entities.AuthResponse, error)
	GetUserByID(ctx context.Context, id string) (*entities.UserResponse, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*entities.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req *entities.UpdateUserRequest) (*entities.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}
