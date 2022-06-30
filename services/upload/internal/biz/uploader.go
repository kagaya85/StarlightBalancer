package biz

import (
	"context"

	v1 "starlight/api/services/process/v1"
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

func (uc *UploaderUsecase) Call(ctx context.Context, selector client.Selector) error {
	ep, err := selector("ProcessService")
	if err != nil {
		return err
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
	client := v1.NewProcessServiceClient(conn)
	reply, err := client.Process(ctx, &v1.ProcessRequest{Id: "2233"})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("[grpc] Process reply %+v\n", reply)
	return nil
}
