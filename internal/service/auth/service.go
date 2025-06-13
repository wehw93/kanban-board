package auth

type Service struct {
	JWTSecret string
}

func NewService(JWTSecret string) *Service {
	return &Service{
		JWTSecret: JWTSecret,
	}
}

func (s *Service) GetJWTSecret() string {
	return s.JWTSecret
}
