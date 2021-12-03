package service

import (
	pb "chat/api/group/service/v1"
	"chat/app/group/service/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GroupService struct {
	pb.UnimplementedGroupServer

	up  *biz.GroupRepository
	log *zap.Logger
}

func NewUserService(up *biz.GroupRepository, logger *zap.Logger) *GroupService {
	return &GroupService{up: up, log: logger}
}

func (s *GroupService) GetGroupInfo(ctx context.Context, req *pb.GetGroupInfoRequest) (*pb.GetGroupInfoResponse, error) {
	s.log.Info("receive msg")
	group := &biz.Group{}
	err := s.up.GetGroupInfo(ctx, group)
	if err != nil {
		return nil, errors.New(400, "zhaobudao", "choisdfasd")
	}
	return &pb.GetGroupInfoResponse{
		Id:        1,
		GroupName: "test group",
	}, nil
}
