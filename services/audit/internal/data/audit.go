package data

import (
	"context"

	"starlight/services/audit/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type auditRepo struct {
	data *Data
	log  *log.Helper
}

// NewAuditRepo .
func NewAuditRepo(data *Data, logger log.Logger) biz.AuditRepo {
	return &auditRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *auditRepo) Save(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *auditRepo) Update(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *auditRepo) FindByID(context.Context, int64) (*biz.Item, error) {
	return nil, nil
}

func (r *auditRepo) ListByHello(context.Context, string) ([]*biz.Item, error) {
	return nil, nil
}

func (r *auditRepo) ListAll(context.Context) ([]*biz.Item, error) {
	return nil, nil
}
