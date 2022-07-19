package data

import (
	"context"

	"starlight/services/timeline/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type timelineRepo struct {
	data *Data
	log  *log.Helper
}

// NewtimelineRepo .
func NewTimelineRepo(data *Data, logger log.Logger) biz.TimelineRepo {
	return &timelineRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *timelineRepo) Save(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *timelineRepo) Update(ctx context.Context, g *biz.Item) (*biz.Item, error) {
	return g, nil
}

func (r *timelineRepo) FindByID(context.Context, int64) (*biz.Item, error) {
	return nil, nil
}

func (r *timelineRepo) ListByHello(context.Context, string) ([]*biz.Item, error) {
	return nil, nil
}

func (r *timelineRepo) ListAll(context.Context) ([]*biz.Item, error) {
	return nil, nil
}
