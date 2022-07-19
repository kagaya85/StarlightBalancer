package data

import (
	"context"

	"starlight/services/push/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type pusherRepo struct {
	data *Data
	log  *log.Helper
}

// NewPusherRepo .
func NewPusherRepo(data *Data, logger log.Logger) biz.PusherRepo {
	return &pusherRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *pusherRepo) Save(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *pusherRepo) Update(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *pusherRepo) FindByID(context.Context, int64) (*biz.Item, error) {
	return nil, nil
}

func (r *pusherRepo) ListByHello(context.Context, string) ([]*biz.Item, error) {
	return nil, nil
}

func (r *pusherRepo) ListAll(context.Context) ([]*biz.Item, error) {
	return nil, nil
}
