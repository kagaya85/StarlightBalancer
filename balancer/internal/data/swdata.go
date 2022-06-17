package data

import (
	"context"
	"encoding/base64"
	assets "starlight/balancer/asserts"
	"starlight/balancer/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/machinebox/graphql"
	api "skywalking.apache.org/repo/goapi/query"
)

type TraceData struct {
	conf   *conf.Data_Skywalking
	client graphql.Client
	log    *log.Helper
}

func NewTraceData(c *conf.Data, logger log.Logger) (*TraceData, error) {
	return &TraceData{
		conf:   c.Skywalking,
		client: *graphql.NewClient(c.Skywalking.GetAddress()),
		log:    log.NewHelper(logger),
	}, nil
}

func (t *TraceData) Execute(ctx context.Context, req *graphql.Request, resp interface{}) error {
	t.setAuthorization(req)
	return t.client.Run(ctx, req, resp)
}

func (t *TraceData) setAuthorization(req *graphql.Request) {
	username := t.conf.GetUsername()
	password := t.conf.GetPassword()
	authorization := ""

	if authorization == "" && username != "" && password != "" {
		authorization = "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	}

	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
}

func (t *TraceData) QueryBasicTraces(ctx context.Context, condition *api.TraceQueryCondition) (api.TraceBrief, error) {
	var rsp map[string]api.TraceBrief

	req := graphql.NewRequest(assets.Read("graphql/Traces.graphql"))
	req.Var("condition", condition)
	err := t.Execute(ctx, req, &rsp)

	return rsp["result"], err
}

func (t *TraceData) QueryTraceIDs(ctx context.Context, startTime time.Time, endTime time.Time, queryState api.TraceState) []string {
	// TODO 采样
	swTimeLayout := "2006-01-02 150405"
	interval := 30 * time.Minute
	itv := interval
	pageNum := 1

	traceIDs := []string{}

	for startTime.Before(endTime) && startTime.Before(time.Now()) {
		start := startTime.Format(swTimeLayout)
		end := startTime.Add(itv).Format(swTimeLayout)

		condition := &api.TraceQueryCondition{
			QueryDuration: &api.Duration{
				Start: start,
				End:   end,
				Step:  api.StepSecond,
			},
			TraceState: queryState,
			QueryOrder: api.QueryOrderByDuration,
			Paging: &api.Pagination{
				PageNum:  &pageNum,
				PageSize: 10000,
			},
		}

		traceBrief, err := t.QueryBasicTraces(ctx, condition)
		if err != nil {
			itv /= 2
			log.Info("get traceID faild, try to use interval of %vs", itv.Seconds())
			if itv <= 0 {
				log.Info("query trace id faild: %s", err)
				itv = 1 * time.Minute
				startTime = startTime.Add(itv)
			}
			continue
		}

		if len(traceBrief.Traces) > 0 {
			log.Info("get %d traceIDs from %q to %q", len(traceBrief.Traces), start, end)
		}

		for _, trace := range traceBrief.Traces {
			traceIDs = append(traceIDs, trace.TraceIds...)
		}

		startTime = startTime.Add(itv)

		if itv < interval {
			itv *= 2
		}
	}

	return traceIDs
}

func (t *TraceData) QueryTraces(ctx context.Context, start time.Time, end time.Time) []api.Trace {
	traceIDs := t.QueryTraceIDs(ctx, start, end, api.TraceStateAll)
	traces := make([]api.Trace, len(traceIDs))
	req := graphql.NewRequest(assets.Read("graphql/Trace.graphql"))

	for _, traceID := range traceIDs {
		var rsp map[string]api.Trace

		req.Var("traceId", traceID)
		if err := t.Execute(ctx, req, &rsp); err != nil {
			log.Error("graphql execute error: %s", err)
			continue
		}

		traces = append(traces, rsp["result"])
	}
	return traces
}
