package auth

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type RegisterResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Token   string `json:"token"`
}
