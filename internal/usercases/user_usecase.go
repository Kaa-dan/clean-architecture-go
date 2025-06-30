package usercases

import (
	"context"
	"time"

	"github.com/kaa-dan/clean-architecture-go/internal/domain/entities"
	"github.com/kaa-dan/clean-architecture-go/internal/domain/services"
	"github.com/kaa-dan/clean-architecture-go/internal/infrastructure/repositories"
	"github.com/kaa-dan/clean-architecture-go/internal/infrastructure/security"
	"github.com/kaa-dan/clean-architecture-go/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userUseCase struct {
	userRepo        repositories.UserRepository
	jwtManager      *security.JWTManager
	passwordManager *security.PasswordManager
}

func NewUserUseCase(
	userRepo repositories.UserRepository,
	jwtManager *security.JWTManager,
	passwordManager *security.PasswordManager,
) services.UserService {
	return &userUseCase{
		userRepo:        userRepo,
		jwtManager:      jwtManager,
		passwordManager: passwordManager,
	}
}

func (u *userUseCase) SignUp(ctx context.Context, req *entities.SignUpRequest) (*entities.AuthResponse, error) {
	// Check if user already exists
	if _, err := u.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, errors.ErrUserAlreadyExists
	}

	if _, err := u.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, errors.ErrUsernameAlreadyExists
	}

	// Hash password
	hashedPassword, err := u.passwordManager.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entities.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		Role:      string(entities.RoleUser),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := u.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &entities.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (u *userUseCase) SignIn(ctx context.Context, req *entities.SignInRequest) (*entities.AuthResponse, error) {
	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Verify password
	if err := u.passwordManager.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.ErrUserInactive
	}

	// Generate token
	token, err := u.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &entities.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (u *userUseCase) GetUserByID(ctx context.Context, id string) (*entities.UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.ErrInvalidUserID
	}

	user, err := u.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (u *userUseCase) GetAllUsers(ctx context.Context, limit, offset int) ([]*entities.UserResponse, error) {
	users, err := u.userRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []*entities.UserResponse
	for _, user := range users {
		response := user.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

func (u *userUseCase) UpdateUser(ctx context.Context, id string, req *entities.UpdateUserRequest) (*entities.UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.ErrInvalidUserID
	}

	// Get existing user
	user, err := u.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Username != nil {
		// Check if username is already taken by another user
		if existingUser, err := u.userRepo.GetByUsername(ctx, *req.Username); err == nil && existingUser.ID != objectID {
			return nil, errors.ErrUsernameAlreadyExists
		}
		user.Username = *req.Username
	}

	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(ctx, objectID, user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (u *userUseCase) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.ErrInvalidUserID
	}

	return u.userRepo.Delete(ctx, objectID)
}
