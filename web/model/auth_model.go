package model

type Auth struct {
	UID string `json:"uid"` // user unique id
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}

type OtpRequest struct {
	CountryCode string `json:"countryCode" validate:"required"`
	Phone       string `json:"phone" validate:"required,max=15"`
}

type VerifyOtpRequest struct {
	DeviceID      string `json:"deviceId" validate:"required"`
	Otp           string `json:"otp" validate:"required"`
	SignatureCode string `json:"signatureCode" validate:"required"`
}
type VerifyOtpResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LogoutRequest struct {
	DeviceID     string `json:"deviceId" validate:"required"`
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type LogoutResponse struct {
	IsLoggedOut bool `json:"isLoggedOut"`
}

type RefreshTokenRequest struct {
	DeviceID     string `json:"deviceId" validate:"required"`
	RefreshToken string `json:"refreshToken" validate:"required"`
}
