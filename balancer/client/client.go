package client

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"

	v1 "starlight/api/balancer/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var ErrNoEndpoint = errors.New("no endpoint")

type Selector func(service string) (endpoint string, err error)

type Weight struct {
	endpoint string
	weight   int
}

type BalancerClient struct {
	l    sync.RWMutex
	list map[string][]Weight // key: service, value: weight list

	serverAddr  string
	maxRetry    int
	serviceName string
	port        string

	log log.Helper
}

func NewBalancerClient(serverAddr string, maxRetry int, serviceName string, port string, logger log.Logger) *BalancerClient {
	return &BalancerClient{
		serverAddr:  serverAddr,
		maxRetry:    maxRetry,
		serviceName: serviceName,
		port:        port,
		log:         *log.NewHelper(logger),
	}
}

func (b *BalancerClient) Sync(ctx context.Context, instanceId string) error {
	errCount := 0
	for {
		conn, err := grpc.DialInsecure(ctx, grpc.WithEndpoint(b.serverAddr))
		if err != nil {
			if errCount < b.maxRetry {
				errCount++
				log.Errorf("connect %s error: %s", b.serverAddr, err)
				<-time.After(3 * time.Second)
				continue
			}
			return err
		}

		errCount = 0
		client := v1.NewWeightUpdaterClient(conn)
		svcInfo, upstreamInfo, err := listServiceInfo(b.serviceName, b.port)
		if err != nil {
			return err
		}

		ip := os.Getenv("POD_IP")
		if ip == "" {
			ip = "127.0.0.1"
		}
		req := &v1.UpdateRequeset{
			Instance: instanceId,
			PodIP:    ip,
			Pod:      os.Getenv("POD_NAME"),
			Node:     os.Getenv("NODE_NAME"),
			Zone:     os.Getenv("ZONE_NAME"),
			Info:     svcInfo,
			Upstream: upstreamInfo,
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

func (b *BalancerClient) sync(stream v1.WeightUpdater_UpdateClient) error {
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
		b.l.Lock()
		b.list = list
		b.l.Unlock()
		b.log.Infof("update weight list: %v", list)
	}
	return nil
}

func (b *BalancerClient) Default(service string) (string, error) {
	return b.Random(service)
}

func (b *BalancerClient) Random(service string) (string, error) {
	b.l.RLock()
	defer b.l.RUnlock()
	list := b.list[service]
	n := len(list)
	if n == 0 {
		return "", ErrNoEndpoint
	}
	return list[rand.Intn(n)].endpoint, nil
}

func (b *BalancerClient) WRandom(service string) (string, error) {
	b.l.RLock()
	defer b.l.RUnlock()
	list := b.list[service]
	n := len(list)
	if n == 0 {
		return "", ErrNoEndpoint
	}
	return list[rand.Intn(n)].endpoint, nil
}

func (b *BalancerClient) DWRandom(service string) (string, error) {
	b.l.RLock()
	defer b.l.RUnlock()
	list := b.list[service]
	n := len(list)
	if n == 0 {
		return "", ErrNoEndpoint
	}
	return list[rand.Intn(n)].endpoint, nil
}

func (b *BalancerClient) RR(service string) (string, error) {
	b.l.RLock()
	defer b.l.RUnlock()
	list := b.list[service]
	n := len(list)
	if n == 0 {
		return "", ErrNoEndpoint
	}
	return list[rand.Intn(n)].endpoint, nil
}

func (b *BalancerClient) WRR(service string) (string, error) {
	b.l.RLock()
	defer b.l.RUnlock()
	list := b.list[service]
	n := len(list)
	if n == 0 {
		return "", ErrNoEndpoint
	}
	return list[rand.Intn(n)].endpoint, nil
}

func (b *BalancerClient) DWRR(service string) (string, error) {
	b.l.RLock()
	defer b.l.RUnlock()
	list := b.list[service]
	n := len(list)
	if n == 0 {
		return "", ErrNoEndpoint
	}
	return list[rand.Intn(n)].endpoint, nil
}

func listServiceInfo(serviceName, port string) (services []*v1.ServiceInfo, upstream []*v1.ServiceInfo, err error) {
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Services().Len(); i++ {
			var (
				methods []string
				sd      = fd.Services().Get(i)
			)
			if sd.Name() == "ServerReflection" || sd.Name() == "WeightUpdater" || sd.Name() == "Metadata" || sd.Name() == "Health" {
				continue
			}
			for j := 0; j < sd.Methods().Len(); j++ {
				md := sd.Methods().Get(j)
				mName := string(md.Name())
				methods = append(methods, mName)
			}
			if sd.Name() == protoreflect.Name(serviceName) {
				services = append(services, &v1.ServiceInfo{
					Service:    string(sd.Name()),
					Port:       port,
					Operations: methods,
				})
			} else {
				upstream = append(upstream, &v1.ServiceInfo{
					Service:    string(sd.Name()),
					Port:       "-1",
					Operations: methods,
				})
			}
		}
		return true
	})
	return
}
