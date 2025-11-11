package repository

import (
	"backend/core/repo"
	"backend/web/entity"

	"github.com/sirupsen/logrus"
)

type BranchsRepository struct {
	repo.Repository[entity.Branch]
	Log *logrus.Logger
}

func NewBranchsRepository(log *logrus.Logger) *BranchsRepository {
	return &BranchsRepository{
		Log: log,
	}
}
