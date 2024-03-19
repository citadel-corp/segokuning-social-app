package user

type UserRegisterResponse struct {
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Name        string  `json:"name"`
	AccessToken string  `json:"accessToken"`
}

type UserLoginResponse struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}
