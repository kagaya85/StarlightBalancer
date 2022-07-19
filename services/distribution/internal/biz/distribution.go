package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Item is a Item model.
type Item struct {
	Hello string
}

// DistributionRepo is a Greater repo.
type DistributionRepo interface {
	Save(context.Context, *Item) (*Item, error)
	Update(context.Context, *Item) (*Item, error)
	FindByID(context.Context, int64) (*Item, error)
	ListByHello(context.Context, string) ([]*Item, error)
	ListAll(context.Context) ([]*Item, error)
}

// DistributionUsecase is a Greeter usecase.
type DistributionUsecase struct {
	repo DistributionRepo
	log  *log.Helper
}

// NewDistributionUsecase new a Greeter usecase.
func NewDistributionUsecase(repo DistributionRepo, logger log.Logger) *DistributionUsecase {
	return &DistributionUsecase{repo: repo, log: log.NewHelper(logger)}
}
