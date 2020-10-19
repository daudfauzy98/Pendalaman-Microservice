package entity

type AuthResponse struct {
	Code         int    `json:"code"`
	Status       string `json:"status, omitemty"`
	ErrorDetails string `json:"error_details, omitemty"`
	ErrorType    string `json:"error_type, omitemty"`
	Data         Data
}

type Data struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
