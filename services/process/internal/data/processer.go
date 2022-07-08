package data

import (
	"context"

	"starlight/services/process/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type processerRepo struct {
	data *Data
	log  *log.Helper
}

// NewProcesserRepo .
func NewProcesserRepo(data *Data, logger log.Logger) biz.ProcesserRepo {
	return &processerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *processerRepo) Save(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *processerRepo) Update(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *processerRepo) FindByID(context.Context, int64) (*biz.Item, error) {
	return nil, nil
}

func (r *processerRepo) ListByHello(context.Context, string) ([]*biz.Item, error) {
	return nil, nil
}

func (r *processerRepo) ListAll(context.Context) ([]*biz.Item, error) {
	return nil, nil
}
