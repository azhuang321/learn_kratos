package data

import (
	"chat/app/group/service/internal/biz"
	"context"
	"go.uber.org/zap"
)

type groupRepo struct {
	data *Data
	log  *zap.Logger
}

func NewGroupRepo(data *Data, logger *zap.Logger) biz.GroupRepo {
	return &groupRepo{
		data: data,
		log:  logger,
	}
}

func (g groupRepo) GetGroupInfo(ctx context.Context, user *biz.Group) error {
	return nil
}
