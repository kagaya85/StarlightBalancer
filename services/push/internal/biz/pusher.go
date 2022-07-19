package biz

import (
	"context"

	v1 "starlight/api/services/timeline/v1"
	"starlight/balancer/client"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// Item is a Item model.
type Item struct {
	Hello string
}

// PusherRepo is a Greater repo.
type PusherRepo interface {
	Save(context.Context, *Item) (*Item, error)
	Update(context.Context, *Item) (*Item, error)
	FindByID(context.Context, int64) (*Item, error)
	ListByHello(context.Context, string) ([]*Item, error)
	ListAll(context.Context) ([]*Item, error)
}

// PusherUsecase is a Greeter usecase.
type PusherUsecase struct {
	repo PusherRepo
	log  *log.Helper
}

// NewPusherUsecase new a Greeter usecase.
func NewPusherUsecase(repo PusherRepo, logger log.Logger) *PusherUsecase {
	return &PusherUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *PusherUsecase) Call(ctx context.Context, selector client.Selector) error {
	ep, err := selector("TimelineService")
	if err != nil {
		log.Errorf("selector error %+v\n", err)
	}
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(ep),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := v1.NewTimelineServiceClient(conn)
	reply, err := client.PushTimeline(ctx, &v1.PushTimelineRequest{UserId: "2233", Items: []string{"hello", "world"}})
	if err != nil {
		log.Error(err)
	}
	log.Infof("[grpc] timeline service reply %+v\n", reply)
	return nil
}
