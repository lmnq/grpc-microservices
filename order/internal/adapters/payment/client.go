package payment

import (
	"context"
	"log"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	paymentpb "github.com/lmnq/microservices-proto/generated/golang/payment"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	payment paymentpb.PaymentClient
}

func NewAdapter(paymentServiceURL string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithChainUnaryInterceptor(
		timeout.UnaryClientInterceptor(time.Second*5),
		circuitBreakerClientUnaryInterceptor(),
		retryClientUnaryInterceptor(),
	))
	conn, err := grpc.NewClient(paymentServiceURL, opts...)
	if err != nil {
		return nil, err
	}
	client := paymentpb.NewPaymentClient(conn)
	return &Adapter{payment: client}, nil
}

func retryClientUnaryInterceptor() grpc.UnaryClientInterceptor {
	return retry.UnaryClientInterceptor(
		retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
		retry.WithMax(3),
		retry.WithBackoff(retry.BackoffLinearWithJitter(time.Second, 0.3)),
	)
}

func circuitBreakerClientUnaryInterceptor() grpc.UnaryClientInterceptor {
	cb := gobreaker.NewCircuitBreaker[interface{}](gobreaker.Settings{
		Name:        "payment",
		MaxRequests: 3,
		Interval:    5 * time.Second,
		Timeout:     3 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio > 0.5
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("CircuitBreaker %s state changed from %s to %s", name, from, to)
		},
	})
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		_, cbErr := cb.Execute(func() (interface{}, error) {
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				return nil, err
			}
			return nil, nil
		})
		return cbErr
	}
}
