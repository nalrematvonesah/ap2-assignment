package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func LoggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	start := time.Now()

	res, err := handler(ctx, req)

	duration := time.Since(start)

	log.Printf("[GRPC][UNARY] %s | %v | error: %v",
		info.FullMethod,
		duration,
		err,
	)

	return res, err
}

func LoggingStreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {

	start := time.Now()

	err := handler(srv, ss)

	duration := time.Since(start)

	log.Printf("[GRPC][STREAM] %s | %v | error: %v",
		info.FullMethod,
		duration,
		err,
	)

	return err
}
