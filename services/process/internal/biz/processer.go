package biz

import (
	"context"

	v1 "starlight/api/services/audit/v1"
	"starlight/balancer/client"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// Item is a Item model.
type Item struct {
	Hello string
}

// ProcesserRepo is a Greater repo.
type ProcesserRepo interface {
	Save(context.Context, *Item) (*Item, error)
	Update(context.Context, *Item) (*Item, error)
	FindByID(context.Context, int64) (*Item, error)
	ListByHello(context.Context, string) ([]*Item, error)
	ListAll(context.Context) ([]*Item, error)
}

// ProcesserUsecase is a Processer usecase.
type ProcesserUsecase struct {
	repo ProcesserRepo
	log  *log.Helper
}

// NewProcesserUsecase new a Processer usecase.
func NewProcesserUsecase(repo ProcesserRepo, logger log.Logger) *ProcesserUsecase {
	return &ProcesserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *ProcesserUsecase) Call(ctx context.Context, selector client.Selector) string {
	ep, err := selector("AuditService")
	if err != nil {
		log.Errorf("selector error %+v\n", err)
		return ""
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
	client := v1.NewAuditServiceClient(conn)
	reply, err := client.Audit(ctx, &v1.AuditRequest{Id: "2233"})
	if err != nil {
		log.Error(err)
	}
	log.Infof("[grpc] Audit reply %+v\n", reply)
	return reply.GetResult()
}
