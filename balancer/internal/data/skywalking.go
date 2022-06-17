package data

import (
	"context"
	"starlight/balancer/internal/biz"
	"strconv"
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
		for _, span := range trace.Spans {
			st := time.UnixMilli(span.StartTime)
			et := time.UnixMilli(span.EndTime)
			spans = append(spans, biz.Span{
				SpanID:       strconv.Itoa(span.SpanID),
				ParentSpanID: strconv.Itoa(span.ParentSpanID),
				TraceID:      span.TraceID,

				Start:    st,
				Duration: et.Sub(st),

				Service:   span.ServiceCode,
				Instance:  span.ServiceInstanceName,
				Operation: *span.EndpointName,
			})
		}
	}
	return spans
}
