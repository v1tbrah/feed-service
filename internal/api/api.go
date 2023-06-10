package api

import (
	"context"
	"fmt"
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/v1tbrah/feed-service/config"
	"github.com/v1tbrah/feed-service/fpbapi"
)

type API struct {
	server *grpc.Server
	cache  Cache
	fpbapi.UnimplementedFeedServiceServer
}

func New(cache Cache) (newAPI *API) {
	newAPI = &API{
		server: grpc.NewServer(grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				interceptorLog,
			),
		)),
		cache: cache,
	}

	fpbapi.RegisterFeedServiceServer(newAPI.server, newAPI)

	return newAPI
}

func (a *API) StartServing(ctx context.Context, cfg config.GRPC, shutdownSig <-chan os.Signal) (err error) {
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	listen, errListen := net.Listen("tcp", addr)
	if errListen != nil {
		return fmt.Errorf("net listen tcp %s server: %w", addr, errListen)
	}

	serveEndSig := make(chan struct{})

	go func() {
		log.Info().Str("addr", addr).Msg("starting gRPC server")
		if err = a.server.Serve(listen); err != nil {
			err = fmt.Errorf("serve %s server: %w", addr, err)
		}
		serveEndSig <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-shutdownSig:
		return err
	case <-serveEndSig:
		return err
	}
}

func (a *API) GracefulStop(ctx context.Context) (err error) {
	gracefulStopEnded := make(chan struct{})

	go func() {
		a.server.GracefulStop()
		gracefulStopEnded <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-gracefulStopEnded:
		return err
	}
}
