package api

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/v1tbrah/feed-service/config"
	"github.com/v1tbrah/feed-service/fpbapi"
	"github.com/v1tbrah/promcli"
)

type API struct {
	server *grpc.Server

	httpServer *http.Server

	promCli *promcli.HTTPReg

	cache Cache

	fpbapi.UnimplementedFeedServiceServer
}

func New(httpCfg config.HTTP, cache Cache) (newAPI *API) {
	newAPI = &API{
		cache:   cache,
		promCli: promcli.NewHTTP("feed_service", "api"),
	}

	newAPI.server = grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			newAPI.interceptor,
		),
	))

	fpbapi.RegisterFeedServiceServer(newAPI.server, newAPI)

	newAPI.httpServer = &http.Server{Addr: net.JoinHostPort(httpCfg.Host, httpCfg.Port)}
	httpRouter := chi.NewRouter()
	httpRouter.Handle("/metrics", promhttp.Handler())
	newAPI.httpServer.Handler = httpRouter

	return newAPI
}

func (a *API) StartServing(ctx context.Context, cfgGRPC config.GRPC, shutdownSig <-chan os.Signal) (err error) {
	grpcAddr := net.JoinHostPort(cfgGRPC.Host, cfgGRPC.Port)
	listen, errListen := net.Listen("tcp", grpcAddr)
	if errListen != nil {
		return errors.Wrapf(errListen, "net listen tcp %s server", grpcAddr)
	}

	grpcServeEndSig := make(chan struct{})
	go func() {
		log.Info().Str("addr", grpcAddr).Msg("starting gRPC server")
		if err = a.server.Serve(listen); err != nil {
			err = errors.Wrapf(err, "serve grpc %s server", grpcAddr)
		}
		grpcServeEndSig <- struct{}{}
	}()

	httpServeEndSig := make(chan struct{})
	go func() {
		log.Info().Str("addr", a.httpServer.Addr).Msg("starting HTTP server")
		if err = a.httpServer.ListenAndServe(); err != nil {
			err = errors.Wrapf(err, "serve http %s server", a.httpServer.Addr)
		}
		httpServeEndSig <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-shutdownSig:
		return err
	case <-grpcServeEndSig:
		return err
	case <-httpServeEndSig:
		return err
	}
}

func (a *API) GracefulStop(ctx context.Context) (err error) {
	gracefulStopEnded := make(chan struct{})

	go func() {
		if err = a.httpServer.Shutdown(ctx); err != nil {
			err = errors.Wrap(err, "http server shutdown")
		}

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
