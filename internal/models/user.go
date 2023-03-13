package models

import "time"

type User struct {
	ID          string    `json:"id" bson:"_id"`
	Email       string    `json:"email,omitempty" bson:"email,omitempty"`
	DiscordID   string    `json:"discord_id,omitempty" bson:"discord_id,omitempty"`
	OtpEnabled  bool      `json:"otp_enabled" bson:"otp_enabled"`
	OtpVerified bool      `json:"otp_verified" bson:"otp_verified"`
	OtpSecret   string    `json:"otp_secret" bson:"otp_secret"`
	OtpUrl      string    `json:"otp_url" bson:"otp_url"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type UserPayload struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	DiscordID string `json:"discord_id"`
}

type UserResponse struct {
	CommonResponse
	User struct {
		ID          string    `json:"id"`
		Email       string    `json:"email"`
		DiscordID   string    `json:"discord_id"`
		OtpEnabled  bool      `json:"otp_enabled"`
		OtpVerified bool      `json:"otp_verified"`
		CreatedAt   time.Time `json:"created_at"`
	} `json:"user"`
}

type OtpPayload struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type OtpResponse struct {
	CommonResponse
	Otp struct {
		OtpEnabled  bool   `json:"otp_enabled,omitempty"`
		OtpVerified bool   `json:"otp_verified,omitempty"`
		OtpValid    bool   `json:"otp_valid,omitempty"`
		OtpSecret   string `json:"otp_secret,omitempty"`
		OtpUrl      string `json:"otp_url,omitempty"`
	} `json:"otp"`
}
