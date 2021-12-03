package service

import (
	pb "chat/api/user/service/v1"
	"chat/app/user/service/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	pb.UnimplementedUserServer

	up  *biz.UserRepository
	log *zap.Logger
}

func NewUserService(up *biz.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{up: up, log: logger}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*emptypb.Empty, error) {
	user := &biz.User{}
	user.Username = req.GetUsername()
	user.Password = req.GetPassword()
	err := s.up.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.New(400, "zhaobudao", "choisdfasd")
		//return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) UserGroupInfo(ctx context.Context, req *pb.UserGroupInfoRequest) (*pb.UserGroupInfoResponse, error) {
	s.up.GetUserGroupInfo(ctx, &biz.User{})

	return &pb.UserGroupInfoResponse{
		GroupId:   1,
		GroupName: "test group",
	}, nil
}
