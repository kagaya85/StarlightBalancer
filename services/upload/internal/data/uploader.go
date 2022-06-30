package data

import (
	"context"

	"starlight/services/upload/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type uploaderRepo struct {
	data *Data
	log  *log.Helper
}

// NewUploaderRepo .
func NewUploaderRepo(data *Data, logger log.Logger) biz.UploaderRepo {
	return &uploaderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *uploaderRepo) Save(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *uploaderRepo) Update(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *uploaderRepo) FindByID(context.Context, int64) (*biz.Item, error) {
	return nil, nil
}

func (r *uploaderRepo) ListByHello(context.Context, string) ([]*biz.Item, error) {
	return nil, nil
}

func (r *uploaderRepo) ListAll(context.Context) ([]*biz.Item, error) {
	return nil, nil
}
