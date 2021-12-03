package biz

import (
	"context"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	GroupInfo(ctx context.Context, user *User) error
}

type UserRepository struct {
	repo UserRepo
	log  *zap.Logger
}

func NewUserRepository(repo UserRepo, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		repo: repo,
		log:  logger,
	}
}

// GeneratePassword 加密用户密码
func generatePassword(userPwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPwd), bcrypt.DefaultCost)
}

func (u *UserRepository) CreateUser(ctx context.Context, user *User) error {
	pwdByte, err := generatePassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(pwdByte)
	return u.repo.Create(ctx, user)
}

func (u *UserRepository) GetUserGroupInfo(ctx context.Context, user *User) error {
	return u.repo.GroupInfo(ctx, user)
}
