package data

import (
	"context"
	"starlight/balancer/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type traceSource struct {
	traceData *TraceData
	logger    *log.Helper
}

func NewTraceSource(traceData *TraceData, logger log.Logger) biz.TraceSource {
	return &traceSource{
		traceData: traceData,
		logger:    log.NewHelper(logger),
	}
}

func (s *traceSource) ListSpanFrom(context context.Context, time time.Time) []biz.Span {
	return []biz.Span{}
}
