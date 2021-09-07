package server

import (
	"context"
	"log"
	"strings"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/pkg/jwtauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Unary returns a server interceptor function to authenticate and authorize unary RPC
func Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		// Check if token is valid and not expired yet
		err := Authorize(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream returns a server interceptor function to authenticate and authorize stream RPC
func Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("--> stream interceptor: ", info.FullMethod)

		// Check if token is valid and not expired yet
		err := Authorize(stream.Context())
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func Authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := strings.TrimSpace(values[0])
	auth := jwtauth.NewJWTAuth()
	_, err := auth.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "token validation failed %v", err)
	}

	return nil
}
