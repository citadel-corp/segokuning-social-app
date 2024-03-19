package user

import (
	"context"
	"fmt"
	"time"

	"github.com/citadel-corp/segokuning-social-app/internal/common/id"
	"github.com/citadel-corp/segokuning-social-app/internal/common/jwt"
	"github.com/citadel-corp/segokuning-social-app/internal/common/password"
)

type Service interface {
	Create(ctx context.Context, req CreateUserPayload) (*UserRegisterResponse, error)
	Login(ctx context.Context, req LoginPayload) (*UserLoginResponse, error)
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
