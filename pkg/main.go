package main

import (
	"context"
	"github.com/softonic/homing-pigeon/pkg/middleware"
	"github.com/softonic/homing-pigeon/proto"
	"golang.org/x/time/rate"
	"k8s.io/klog"
	"os"
	"strconv"
)

type ThrottlingMiddleware struct {
	middleware.UnimplementedMiddleware
	limiter *rate.Limiter
}

func NewThrottlingMiddleware(r rate.Limit, b int) *ThrottlingMiddleware {
	if r == 0 || b == 0 {
		return &ThrottlingMiddleware{
			limiter: nil,
		}
	}
	return &ThrottlingMiddleware{
		limiter: rate.NewLimiter(r, b),
	}
}

func (m *ThrottlingMiddleware) Handle(ctx context.Context, req *proto.Data) (*proto.Data, error) {
	if m.limiter != nil {
		// Wait for permission to proceed with request due to throttling
		if err := m.limiter.Wait(ctx); err != nil {
			return nil, err
		}
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

	// Default values
	defaultLimit := 100.0
	defaultBurst := 100

	// Read limit and burst from environment variables
	limitStr := os.Getenv("THROTTLE_LIMIT")
	burstStr := os.Getenv("THROTTLE_BURST")

	// Convert limit and burst to appropriate types
	limit, err := strconv.ParseFloat(limitStr, 64)
	if err != nil || limit < 0 {
		klog.Warningf("Invalid or missing THROTTLE_LIMIT, using default: %v", defaultLimit)
		limit = defaultLimit
	}

	burst, err := strconv.Atoi(burstStr)
	if err != nil || burst < 0 {
		klog.Warningf("Invalid or missing THROTTLE_BURST, using default: %v", defaultBurst)
		burst = defaultBurst
	}

	// Create limiter with the provided limit and burst
	middleware := NewThrottlingMiddleware(rate.Limit(limit), burst)
	middleware.Listen(middleware)

	klog.Flush()
}
