package user

import (
	"context"
	"fmt"
	"time"

	"github.com/citadel-corp/segokuning-social-app/internal/common/id"
	"github.com/citadel-corp/segokuning-social-app/internal/common/jwt"
	"github.com/citadel-corp/segokuning-social-app/internal/common/password"
	"github.com/citadel-corp/segokuning-social-app/internal/common/response"
)

type Service interface {
	Create(ctx context.Context, req CreateUserPayload) (*UserRegisterResponse, error)
	Login(ctx context.Context, req LoginPayload) (*UserLoginResponse, error)
	LinkEmail(ctx context.Context, req LinkEmailPayload, userID string) error
	LinkPhoneNumber(ctx context.Context, req LinkPhoneNumberPayload, userID string) error
	Update(ctx context.Context, req UpdateUserPayload, userID string) error
	List(ctx context.Context, req ListUserPayload) ([]UserListResponse, *response.Pagination, error)
}

type userService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userService{repository: repository}
}

func (s *userService) Create(ctx context.Context, req CreateUserPayload) (*UserRegisterResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	var user *User
	if req.CredentialType == "email" {
		user = &User{
			ID:             id.GenerateStringID(16),
			Name:           req.Name,
			Email:          &req.CredentialValue,
			HashedPassword: hashedPassword,
		}
		err = s.repository.CreateWithEmail(ctx, user)
	} else {
		user = &User{
			ID:             id.GenerateStringID(16),
			Name:           req.Name,
			PhoneNumber:    &req.CredentialValue,
			HashedPassword: hashedPassword,
		}
		err = s.repository.CreateWithPhoneNumber(ctx, user)
	}
	if err != nil {
		return nil, err
	}
	// create access token with signed jwt
	accessToken, err := jwt.Sign(time.Hour*8, fmt.Sprint(user.ID))
	if err != nil {
		return nil, err
	}
	return &UserRegisterResponse{
		Phone:       user.PhoneNumber,
		Email:       user.Email,
		Name:        req.Name,
		AccessToken: accessToken,
	}, nil
}

func (s *userService) Login(ctx context.Context, req LoginPayload) (*UserLoginResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	var user *User
	if req.CredentialType == "email" {
		user, err = s.repository.GetByEmail(ctx, req.CredentialValue)
	} else {
		user, err = s.repository.GetByPhoneNumber(ctx, req.CredentialValue)
	}
	if err != nil {
		return nil, err
	}
	match, err := password.Matches(req.Password, user.HashedPassword)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrWrongPassword
	}
	// create access token with signed jwt
	accessToken, err := jwt.Sign(time.Hour*8, fmt.Sprint(user.ID))
	if err != nil {
		return nil, err
	}
	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	phone := ""
	if user.PhoneNumber != nil {
		phone = *user.PhoneNumber
	}
	return &UserLoginResponse{
		Email:       email,
		Phone:       phone,
		Name:        user.Name,
		AccessToken: accessToken,
	}, nil
}

// LinkEmail implements Service.
func (s *userService) LinkEmail(ctx context.Context, req LinkEmailPayload, userID string) error {
	err := req.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.Email != nil {
		return ErrUserHasEmail
	}
	user.Email = &req.Email
	return s.repository.Update(ctx, user)
}

// LinkPhoneNumber implements Service.
func (s *userService) LinkPhoneNumber(ctx context.Context, req LinkPhoneNumberPayload, userID string) error {
	err := req.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.PhoneNumber != nil {
		return ErrUserHasPhoneNumber
	}
	user.PhoneNumber = &req.Phone
	return s.repository.Update(ctx, user)
}

// Update implements Service.
func (s *userService) Update(ctx context.Context, req UpdateUserPayload, userID string) error {
	err := req.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	user.ImageURL = &req.ImageURL
	user.Name = req.Name
	return s.repository.Update(ctx, user)
}

func (s *userService) List(ctx context.Context, req ListUserPayload) ([]UserListResponse, *response.Pagination, error) {
	req.WithoutUser = true
	return s.repository.List(ctx, req)
}
