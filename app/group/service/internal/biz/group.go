package biz

import (
	"context"
	"go.uber.org/zap"
)

type Group struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type GroupRepo interface {
	GetGroupInfo(ctx context.Context, user *Group) error
}

type GroupRepository struct {
	repo GroupRepo
	log  *zap.Logger
}

func NewGroupRepository(repo GroupRepo, logger *zap.Logger) *GroupRepository {
	return &GroupRepository{
		repo: repo,
		log:  logger,
	}
}

func (u *GroupRepository) GetGroupInfo(ctx context.Context, group *Group) error {
	return u.repo.GetGroupInfo(ctx, group)
}
