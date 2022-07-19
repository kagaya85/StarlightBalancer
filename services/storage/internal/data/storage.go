package data

import (
	"context"

	"starlight/services/storage/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type storageRepo struct {
	data *Data
	log  *log.Helper
}

// NewStorageRepo .
func NewStorageRepo(data *Data, logger log.Logger) biz.StorageRepo {
	return &storageRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *storageRepo) Save(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *storageRepo) Update(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *storageRepo) FindByID(context.Context, int64) (*biz.Item, error) {
	return nil, nil
}

func (r *storageRepo) ListByHello(context.Context, string) ([]*biz.Item, error) {
	return nil, nil
}

func (r *storageRepo) ListAll(context.Context) ([]*biz.Item, error) {
	return nil, nil
}
