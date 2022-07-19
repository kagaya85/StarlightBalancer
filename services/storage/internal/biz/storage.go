package biz

import (
	"context"

	v1 "starlight/api/services/distribution/v1"
	"starlight/balancer/client"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// Item is a Item model.
type Item struct {
	Hello string
}

// StorageRepo is a Storage repo.
type StorageRepo interface {
	Save(context.Context, *Item) (*Item, error)
	Update(context.Context, *Item) (*Item, error)
	FindByID(context.Context, int64) (*Item, error)
	ListByHello(context.Context, string) ([]*Item, error)
	ListAll(context.Context) ([]*Item, error)
}

// StorageUsecase is a Greeter usecase.
type StorageUsecase struct {
	repo StorageRepo
	log  *log.Helper
}

// NewStorageUsecase new a Greeter usecase.
func NewStorageUsecase(repo StorageRepo, logger log.Logger) *StorageUsecase {
	return &StorageUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *StorageUsecase) Call(ctx context.Context, selector client.Selector) error {
	ep, err := selector("DistributionService")
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
	client := v1.NewDistributionServiceClient(conn)
	reply, err := client.Distribute(ctx, &v1.DistributeRequest{Id: "2233"})
	if err != nil {
		log.Error(err)
	}
	log.Infof("[grpc] distribution service reply %+v\n", reply)
	return nil
}
