package models

type User struct {
	ID          string `json:"id" bson:"id"`
	Email       string `json:"email,omitempty" bson:"email,omitempty"`
	DiscordID   string `json:"discord_id,omitempty" bson:"discord_id,omitempty"`
	OtpEnabled  bool   `json:"otp_enabled" bson:"otp_enabled"`
	OtpVerified bool   `json:"otp_verified" bson:"otp_verified"`
	OtpSecret   string `json:"otp_secret" bson:"otp_secret"`
	OtpUrl      string `json:"otp_url" bson:"otp_url"`
}

type UserPayload struct {
	ID        string `json:"id"`
	Email     string `json:"email,omitempty"`
	DiscordID string `json:"discord_id,omitempty"`
}

type OtpPayload struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}
