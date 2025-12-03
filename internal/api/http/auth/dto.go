package authapi

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type RegisterResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
