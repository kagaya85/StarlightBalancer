package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Item is a Item model.
type Item struct {
	Hello string
}

// TimelineRepo is a Greater repo.
type TimelineRepo interface {
	Save(context.Context, *Item) (*Item, error)
	Update(context.Context, *Item) (*Item, error)
	FindByID(context.Context, int64) (*Item, error)
	ListByHello(context.Context, string) ([]*Item, error)
	ListAll(context.Context) ([]*Item, error)
}

// TimelineUsecase is a Greeter usecase.
type TimelineUsecase struct {
	repo TimelineRepo
	log  *log.Helper
}

// NewTimelineUsecase new a Greeter usecase.
func NewTimelineUsecase(repo TimelineRepo, logger log.Logger) *TimelineUsecase {
	return &TimelineUsecase{repo: repo, log: log.NewHelper(logger)}
}
