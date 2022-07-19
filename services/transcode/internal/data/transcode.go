package data

import (
	"context"

	"starlight/services/transcode/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type transcodeRepo struct {
	data *Data
	log  *log.Helper
}

// NewtranscodeRepo .
func NewtranscodeRepo(data *Data, logger log.Logger) biz.TranscodeRepo {
	return &transcodeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *transcodeRepo) Save(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *transcodeRepo) Update(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *transcodeRepo) FindByID(context.Context, int64) (*biz.Item, error) {
	return nil, nil
}

func (r *transcodeRepo) ListByHello(context.Context, string) ([]*biz.Item, error) {
	return nil, nil
}

func (r *transcodeRepo) ListAll(context.Context) ([]*biz.Item, error) {
	return nil, nil
}
