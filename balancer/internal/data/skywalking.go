package data

import (
	"context"
	"starlight/balancer/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type traceSource struct {
	traceData *TraceData
	log       *log.Helper
}

func NewTraceSource(traceData *TraceData, logger log.Logger) biz.TraceSource {
	return &traceSource{
		traceData: traceData,
		log:       log.NewHelper(logger),
	}
}

func (s *traceSource) ListSpan(ctx context.Context, duration time.Duration) []biz.Span {
	end := time.Now()
	start := end.Add(duration)
	traces := s.traceData.QueryTraces(ctx, start, end)
	s.log.Infof("queried %d traces from %s to %s", len(traces), start.Format(time.Layout), end.Format(time.Layout))

	spans := make([]biz.Span, 10)
	for _, trace := range traces {
		insMap := map[int]string{-1: "root"}
		svcMap := map[int]string{-1: "root-service"}
		opMap := map[int]string{-1: "root-operation"}
		// TODO 获取接口依赖 只保留Entry span
		for _, span := range trace.Spans {
			insMap[span.SpanID] = span.ServiceInstanceName
			svcMap[span.SpanID] = span.ServiceCode
			opMap[span.SpanID] = *span.EndpointName
		}
		for _, span := range trace.Spans {
			st := time.UnixMilli(span.StartTime)
			et := time.UnixMilli(span.EndTime)
			spans = append(spans, biz.Span{
				TraceID:  span.TraceID,
				Start:    st,
				Duration: et.Sub(st),

				CallerIns: insMap[span.ParentSpanID],
				CalleeIns: insMap[span.SpanID],
				CallerSvc: svcMap[span.ParentSpanID],
				CalleeSvc: svcMap[span.SpanID],
				CallerOp:  *span.EndpointName,
				CalleeOp:  *span.EndpointName,
			})
		}
	}
	return spans
}
