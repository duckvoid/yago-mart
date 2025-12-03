package service

type AuthService struct {
	userSvc UserService
}

func NewAuthService(userSvc UserService) *AuthService {
	return &AuthService{userSvc: userSvc}
}

func (s *AuthService) Register(username, password string) error {
	return s.userSvc.Create(username, password)
}
