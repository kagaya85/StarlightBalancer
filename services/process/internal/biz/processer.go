package biz

import (
	"context"

	audit "starlight/api/services/audit/v1"
	storage "starlight/api/services/storage/v1"
	transcode "starlight/api/services/transcode/v1"
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

func (uc *ProcesserUsecase) CallAudit(ctx context.Context, selector client.Selector) (string, error) {
	ep, err := selector("AuditService")
	if err != nil {
		log.Errorf("selector error %+v\n", err)
		return "", err
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
	client := audit.NewAuditServiceClient(conn)
	reply, err := client.Audit(ctx, &audit.AuditRequest{Id: "2233"})
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Infof("[grpc] audit service reply %+v\n", reply)
	return reply.GetResult(), nil
}

func (uc *ProcesserUsecase) CallTranscode(ctx context.Context, selector client.Selector) (string, error) {
	ep, err := selector("TranscodeService")
	if err != nil {
		log.Errorf("selector error %+v\n", err)
		return "", err
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
	client := transcode.NewTranscodeServiceClient(conn)
	reply, err := client.Transcode(ctx, &transcode.TranscodeRequest{Source: "www.kagaya.com/foo/bar.mp4"})
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Infof("[grpc] transcode service reply %+v\n", reply)
	return reply.GetTarget(), nil
}

func (uc *ProcesserUsecase) CallStorage(ctx context.Context, selector client.Selector) (string, error) {
	ep, err := selector("StorageService")
	if err != nil {
		log.Errorf("selector error %+v\n", err)
		return "", err
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
	client := storage.NewStorageServiceClient(conn)
	reply, err := client.Save(ctx, &storage.SaveRequest{Name: "kagaya", Data: "foobar"})
	if err != nil {
		log.Error(err)
		return "", err
	}
	log.Infof("[grpc] Audit reply %+v\n", reply)
	return reply.GetResult(), nil
}
