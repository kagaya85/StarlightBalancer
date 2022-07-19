package biz

import (
	"context"

	process "starlight/api/services/process/v1"
	push "starlight/api/services/push/v1"
	"starlight/balancer/client"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// Item is a Item model.
type Item struct {
	Hello string
}

// UploaderRepo is a Greater repo.
type UploaderRepo interface {
	Save(context.Context, *Item) (*Item, error)
	Update(context.Context, *Item) (*Item, error)
	FindByID(context.Context, int64) (*Item, error)
	ListByHello(context.Context, string) ([]*Item, error)
	ListAll(context.Context) ([]*Item, error)
}

// UploaderUsecase is a Greeter usecase.
type UploaderUsecase struct {
	repo UploaderRepo
	log  *log.Helper
}

// NewUploaderUsecase new a Greeter usecase.
func NewUploaderUsecase(repo UploaderRepo, logger log.Logger) *UploaderUsecase {
	return &UploaderUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UploaderUsecase) CallProcesss(ctx context.Context, selector client.Selector) error {
	ep, err := selector("ProcessService")
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
	client := process.NewProcessServiceClient(conn)
	reply, err := client.Process(ctx, &process.ProcessRequest{Id: "2233"})
	if err != nil {
		log.Error(err)
	}
	log.Infof("[grpc] process service reply %+v\n", reply)
	return nil
}

func (uc *UploaderUsecase) CallPush(ctx context.Context, selector client.Selector) error {
	ep, err := selector("PushService")
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
	client := push.NewPushServiceClient(conn)
	reply, err := client.PushVideo(ctx, &push.PushRequest{Id: "2233"})
	if err != nil {
		log.Error(err)
	}
	log.Infof("[grpc] push service reply %+v\n", reply)
	return nil
}
