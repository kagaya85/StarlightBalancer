package main

import (
	"context"
	"flag"
	"os"
	"strconv"
	"strings"

	"starlight/balancer/client"
	"starlight/services/timeline/internal/conf"
	"starlight/services/timeline/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "TimelineService"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	osname, _ = os.Hostname()

	id = osname + "#" + strconv.Itoa(os.Getpid())
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	service.GlobalBalancer = client.NewBalancerClient(bc.Balancer.Addr, int(bc.Balancer.MaxRetry), Name, strings.Split(bc.Server.Grpc.Addr, ":")[1], bc.Balancer.Method, logger)

	if err := c.Watch("balancer.method", func(key string, value config.Value) {
		service.GlobalBalancer.SetMethod(value.Load().(string))
	}); err != nil {
		panic(err)
	}

	go func() {
		if err := service.GlobalBalancer.Sync(context.TODO(), id); err != nil {
			panic(err)
		}
	}()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
