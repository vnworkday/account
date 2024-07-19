//nolint:nonamedreturns
package port

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"go.uber.org/zap"
)

func LoggingMiddleware(logger *zap.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (response any, err error) {
			defer func(begin time.Time) {
				logger.Info("invoke", zap.Error(err), zap.Duration("took", time.Since(begin)))
			}(time.Now())

			return next(ctx, request)
		}
	}
}

func InstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (response any, err error) {
			defer func(begin time.Time) {
				duration.With("success", strconv.FormatBool(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())

			return next(ctx, request)
		}
	}
}
