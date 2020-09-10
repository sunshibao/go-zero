package serverinterceptors

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/rpcx/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestStreamAuthorizeInterceptor(t *testing.T) {
	tests := []struct {
		name     string
		app      string
		token    string
		strict   bool
		hasError bool
	}{
		{
			name:     "strict=false",
			strict:   false,
			hasError: false,
		},
		{
			name:     "strict=true",
			strict:   true,
			hasError: true,
		},
		{
			name:     "strict=true,with token",
			app:      "foo",
			token:    "bar",
			strict:   true,
			hasError: false,
		},
		{
			name:     "strict=true,with error token",
			app:      "foo",
			token:    "error",
			strict:   true,
			hasError: true,
		},
	}

	r := miniredis.NewMiniRedis()
	assert.Nil(t, r.Start())
	defer r.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := redis.NewRedis(r.Addr(), redis.NodeType)
			if len(test.app) > 0 {
				assert.Nil(t, store.Hset("apps", test.app, test.token))
				defer store.Hdel("apps", test.app)
			}

			authenticator, err := auth.NewAuthenticator(store, "apps", test.strict)
			assert.Nil(t, err)
			interceptor := StreamAuthorizeInterceptor(authenticator)
			md := metadata.New(map[string]string{
				"app":   "foo",
				"token": "bar",
			})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			stream := mockedStream{ctx: ctx}
			err = interceptor(nil, stream, nil, func(srv interface{}, stream grpc.ServerStream) error {
				return nil
			})
			if test.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestUnaryAuthorizeInterceptor(t *testing.T) {
	tests := []struct {
		name     string
		app      string
		token    string
		strict   bool
		hasError bool
	}{
		{
			name:     "strict=false",
			strict:   false,
			hasError: false,
		},
		{
			name:     "strict=true",
			strict:   true,
			hasError: true,
		},
		{
			name:     "strict=true,with token",
			app:      "foo",
			token:    "bar",
			strict:   true,
			hasError: false,
		},
		{
			name:     "strict=true,with error token",
			app:      "foo",
			token:    "error",
			strict:   true,
			hasError: true,
		},
	}

	r := miniredis.NewMiniRedis()
	assert.Nil(t, r.Start())
	defer r.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := redis.NewRedis(r.Addr(), redis.NodeType)
			if len(test.app) > 0 {
				assert.Nil(t, store.Hset("apps", test.app, test.token))
				defer store.Hdel("apps", test.app)
			}

			authenticator, err := auth.NewAuthenticator(store, "apps", test.strict)
			assert.Nil(t, err)
			interceptor := UnaryAuthorizeInterceptor(authenticator)
			md := metadata.New(map[string]string{
				"app":   "foo",
				"token": "bar",
			})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			_, err = interceptor(ctx, nil, nil,
				func(ctx context.Context, req interface{}) (interface{}, error) {
					return nil, nil
				})
			if test.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			if test.strict {
				_, err = interceptor(context.Background(), nil, nil,
					func(ctx context.Context, req interface{}) (interface{}, error) {
						return nil, nil
					})
				assert.NotNil(t, err)

				var md metadata.MD
				ctx := metadata.NewIncomingContext(context.Background(), md)
				_, err = interceptor(ctx, nil, nil,
					func(ctx context.Context, req interface{}) (interface{}, error) {
						return nil, nil
					})
				assert.NotNil(t, err)

				md = metadata.New(map[string]string{
					"app":   "",
					"token": "",
				})
				ctx = metadata.NewIncomingContext(context.Background(), md)
				_, err = interceptor(ctx, nil, nil,
					func(ctx context.Context, req interface{}) (interface{}, error) {
						return nil, nil
					})
				assert.NotNil(t, err)
			}
		})
	}
}

type mockedStream struct {
	ctx context.Context
}

func (m mockedStream) SetHeader(md metadata.MD) error {
	return nil
}

func (m mockedStream) SendHeader(md metadata.MD) error {
	return nil
}

func (m mockedStream) SetTrailer(md metadata.MD) {
}

func (m mockedStream) Context() context.Context {
	return m.ctx
}

func (m mockedStream) SendMsg(v interface{}) error {
	return nil
}

func (m mockedStream) RecvMsg(v interface{}) error {
	return nil
}
