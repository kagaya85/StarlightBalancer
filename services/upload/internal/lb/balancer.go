package lb

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	v1 "starlight/api/balancer/v1"

	"starlight/services/upload/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// global weight list for load balance
var GlobalBalancer Balancer

var ErrNoEndpoint = errors.New("no endpoint")

type Selector func(service string) (endpoint string, err error)

type Weight struct {
	endpoint string
	weight   int
}

type Balancer struct {
	mu   sync.Mutex
	list map[string][]Weight // key: service, value: weight list

	conf *conf.Balancer
	log  log.Helper
}

func (b *Balancer) Sync(ctx context.Context, conf *conf.Balancer, svrconf *conf.Server) error {
	b.conf = conf
	errCount := 0
	for {
		conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(conf.Addr))
		if err != nil {
			if errCount < int(conf.MaxRetry) {
				errCount++
				<-time.After(3 * time.Second)
				continue
			}
			return err
		}

		errCount = 0
		client := v1.NewWeightUpdaterClient(conn)
		svcInfo, err := listServiceInfo(svrconf)
		if err != nil {
			return err
		}

		req := &v1.UpdateRequeset{
			Instance: os.Getenv("POD_IP"),
			Pod:      os.Getenv("POD_NAME"),
			Node:     os.Getenv("NODE_NAME"),
			Zone:     os.Getenv("ZONE_NAME"),
			Info:     svcInfo,
		}
		stream, err := client.Update(ctx, req)
		if ctx.Err() != nil {
			// context exceeded
			break
		}
		if err != nil {
			return err
		}
		if err := b.sync(stream); err != nil {
			b.log.Errorf("update stream error %s", err)
			continue
		} else {
			break
		}
	}
	return nil
}

func (b *Balancer) sync(stream v1.WeightUpdater_UpdateClient) error {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		wl := resp.GetWeightList()
		list := make(map[string][]Weight, len(wl))
		for op, inss := range wl {
			for ep, value := range inss.InstanceWeight {
				list[op] = append(list[op], Weight{
					endpoint: ep,
					weight:   int(value),
				})
			}
		}
		b.mu.Lock()
		b.list = list
		b.mu.Unlock()
	}
	return nil
}

func (b *Balancer) Random(service string) (string, error) {
	return "", nil
}

func (b *Balancer) WRandom(service string) (string, error) {
	return "", nil
}

func (b *Balancer) DWRandom(service string) (string, error) {
	return "", nil
}

func (b *Balancer) RR(service string) (string, error) {
	return "", nil
}

func (b *Balancer) WRR(service string) (string, error) {
	return "", nil
}

func (b *Balancer) DWRR(service string) (string, error) {
	return "", nil
}

func listServiceInfo(svrconf *conf.Server) (services []*v1.ServiceInfo, err error) {
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Services().Len(); i++ {
			var (
				methods []string
				sd      = fd.Services().Get(i)
			)
			for j := 0; j < sd.Methods().Len(); j++ {
				md := sd.Methods().Get(j)
				mName := string(md.Name())
				methods = append(methods, mName)
			}
			services = append(services, &v1.ServiceInfo{
				Service:    string(sd.Name()),
				Port:       strings.Split(svrconf.Grpc.Addr, ":")[1],
				Operations: methods,
			})
		}
		return true
	})
	return
}
