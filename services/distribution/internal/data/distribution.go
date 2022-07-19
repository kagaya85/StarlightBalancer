package data

import (
	"context"

	"starlight/services/distribution/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type distributionRepo struct {
	data *Data
	log  *log.Helper
}

// NewDistribuionRepo .
func NewDistribuionRepo(data *Data, logger log.Logger) biz.DistributionRepo {
	return &distributionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *distributionRepo) Save(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *distributionRepo) Update(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *distributionRepo) FindByID(context.Context, int64) (*biz.Item, error) {
	return nil, nil
}

func (r *distributionRepo) ListByHello(context.Context, string) ([]*biz.Item, error) {
	return nil, nil
}

func (r *distributionRepo) ListAll(context.Context) ([]*biz.Item, error) {
	return nil, nil
}
