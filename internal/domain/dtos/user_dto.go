package dtos

// Requests
type RegisterRequestPayload struct {
	Username    string `json:"username" validate:"required,min=6,max=20,ascii"`
	Email       string `json:"email" validate:"required,email"`
	Phonenumber string `json:"phoneNumber" validate:"numeric,min=7"`
	Password    string `json:"password" validate:"required,alphanum,min=6,max=20"`
}

type LoginRequestPayload struct {
	Username string `json:"username" validate:"required,min=6,max=20,ascii"`
	Password string `json:"password" validate:"required,alphanum,min=6,max=20"`
}

type RefreshTokenRequestPayload struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type ForgotPasswordRequestPayload struct {
	Id              string `json:"id" validate:"required"`
	Password        string `json:"password" validate:"required,min=6,max=20,alphanum"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=6,max=20,alphanum,eqfield=password"`
}

type DeleteAccountRequestPayload struct {
	UserId string `form:"userId" validate:"required,numeric"`
}

// Responses
type LoginResponsePayload struct {
	UserId       string `json:"userId"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Phonenumber  string `json:"phoneNumber"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponsePayload struct {
	AccessToken string `json:"accessToken"`
}
