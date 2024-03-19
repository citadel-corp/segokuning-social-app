package user

import "errors"

var (
	ErrUserNotFound                 = errors.New("user not found")
	ErrWrongPassword                = errors.New("wrong password")
	ErrUserPhoneNumberAlreadyExists = errors.New("user phone number already exists")
	ErrUserEmailAlreadyExists       = errors.New("user email already exists")
	ErrUserHasEmail                 = errors.New("user already has email")
	ErrUserHasPhoneNumber           = errors.New("user already has phone number")
	ErrValidationFailed             = errors.New("validation failed")
	ErrCredentialMustExists         = errors.New("credential must be one of phone/email")
)
