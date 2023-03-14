package models

type OtpPayload struct {
	ID        string `json:"id"`
	DiscordID string `json:"discord_id"`
	Token     string `json:"token"`
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
