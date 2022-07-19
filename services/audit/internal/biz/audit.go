package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Item is a Item model.
type Item struct {
	Hello string
}

// AuditRepo is a Audit repo.
type AuditRepo interface {
	Save(context.Context, *Item) (*Item, error)
	Update(context.Context, *Item) (*Item, error)
	FindByID(context.Context, int64) (*Item, error)
	ListByHello(context.Context, string) ([]*Item, error)
	ListAll(context.Context) ([]*Item, error)
}

// AuditUsecase is a Greeter usecase.
type AuditUsecase struct {
	repo AuditRepo
	log  *log.Helper
}

// NewAuditUsecase new a Greeter usecase.
func NewAuditUsecase(repo AuditRepo, logger log.Logger) *AuditUsecase {
	return &AuditUsecase{repo: repo, log: log.NewHelper(logger)}
}
