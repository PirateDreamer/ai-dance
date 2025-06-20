package application

import "context"

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUserInfo(ctx context.Context) {
}
