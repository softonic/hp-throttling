package main

import (
	"context"
	"github.com/softonic/homing-pigeon/pkg/middleware"
	"github.com/softonic/homing-pigeon/proto"
	"golang.org/x/time/rate"
	"k8s.io/klog"
)

type ThrottlingMiddleware struct {
	middleware.UnimplementedMiddleware
	limiter *rate.Limiter
}

func NewThrottlingMiddleware(r rate.Limit, b int) *ThrottlingMiddleware {
	return &ThrottlingMiddleware{
		limiter: rate.NewLimiter(r, b),
	}
}

func (m *ThrottlingMiddleware) Handle(ctx context.Context, req *proto.Data) (*proto.Data, error) {
	// Wait for permission to proceed with request due to throttling
	if err := m.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	// Do things with the INPUT data
	klog.Infof("Pre-Processing %v", *req)

	// Send data to the next middleware and got the response
	resp, err := m.Next(req)

	// Do things with the OUTPUT data
	klog.Infof("Post-Processing %v", *resp)

	return resp, err
}

func main() {
	klog.InitFlags(nil)

	// Limitar a 100 solicitudes por segundo
	limiter := rate.NewLimiter(rate.Limit(100), 100)
	middleware := &ThrottlingMiddleware{
		limiter: limiter,
	}
	middleware.Listen(middleware)

	klog.Flush()
}
