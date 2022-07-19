package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Item is a Item model.
type Item struct {
	Hello string
}

// transcodeRepo is a Greater repo.
type TranscodeRepo interface {
	Save(context.Context, *Item) (*Item, error)
	Update(context.Context, *Item) (*Item, error)
	FindByID(context.Context, int64) (*Item, error)
	ListByHello(context.Context, string) ([]*Item, error)
	ListAll(context.Context) ([]*Item, error)
}

// transcodeUsecase is a Greeter usecase.
type TranscodeUsecase struct {
	repo TranscodeRepo
	log  *log.Helper
}

// NewtranscodeUsecase new a Greeter usecase.
func NewtranscodeUsecase(repo TranscodeRepo, logger log.Logger) *TranscodeUsecase {
	return &TranscodeUsecase{repo: repo, log: log.NewHelper(logger)}
}
