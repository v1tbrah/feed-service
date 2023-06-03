package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/pet-pr-social-network/feed-service/config"
	"gitlab.com/pet-pr-social-network/feed-service/internal/api"
	"gitlab.com/pet-pr-social-network/feed-service/internal/cache"
	"gitlab.com/pet-pr-social-network/feed-service/internal/eventreader"
	"gitlab.com/pet-pr-social-network/feed-service/internal/postcli"
	"gitlab.com/pet-pr-social-network/feed-service/internal/relationcli"
)

func main() {
	newConfig := config.NewDefaultConfig()
	zerolog.SetGlobalLevel(newConfig.LogLvl)

	if err := newConfig.ParseEnv(); err != nil {
		log.Fatal().Err(err).Msg("config.ParseEnv")
	}
	zerolog.SetGlobalLevel(newConfig.LogLvl)

	ctxStart, ctxStartCancel := context.WithCancel(context.Background())

	newCache, err := cache.Init(newConfig.CacheConfig)
	if err != nil {
		log.Fatal().Err(err).Str("config", fmt.Sprintf("%+v", newConfig.CacheConfig)).Msg("storage.Init")
	} else {
		log.Info().Msg("storage initialized")
	}

	newRelationCli, err := relationcli.New(newConfig.RelationCliConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("relationcli.New")
	} else {
		log.Info().Msg("relationcli initialized")
	}

	newPostCli, err := postcli.New(newConfig.PostCliConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("postcli.New")
	} else {
		log.Info().Msg("postcli initialized")
	}

	newMsgReader, err := eventreader.Init(ctxStart, newCache, newRelationCli, newPostCli, newConfig.KafkaConfig)
	if err != nil {
		log.Fatal().Err(err).Interface("config", newConfig.KafkaConfig).Msg("msgrdr.New")
	} else {
		log.Info().Msg("message reader initialized")
	}

	newAPI := api.New(newCache)

	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	errServingCh := make(chan error)
	go func() {
		errServing := newAPI.StartServing(context.Background(), newConfig.GRPCConfig, shutdownSig)
		errServingCh <- errServing
	}()

	select {
	case <-shutdownSig:
		close(shutdownSig)
	case errServing := <-errServingCh:
		if errServing != nil {
			log.Error().Err(errServing).Msg("newAPI.StartServing")
		}
	}

	ctxStartCancel()

	ctxClose, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err = newAPI.GracefulStop(ctxClose); err != nil {
		log.Error().Err(err).Msg("gRPC server graceful stop")
		if err == context.DeadlineExceeded {
			return
		}
	} else {
		log.Info().Msg("gRPC server gracefully stopped")
	}

	if err = newMsgReader.Close(ctxClose); err != nil {
		log.Error().Err(err).Msg("message reader close")
		if err == context.DeadlineExceeded {
			return
		}
	} else {
		log.Info().Msg("message reader closed")
	}

	if err = newRelationCli.Close(ctxClose); err != nil {
		log.Error().Err(err).Msg("relationcli close")
		if err == context.DeadlineExceeded {
			return
		}
	} else {
		log.Info().Msg("relationcli closed")
	}

	if err = newPostCli.Close(ctxClose); err != nil {
		log.Error().Err(err).Msg("postcli close")
		if err == context.DeadlineExceeded {
			return
		}
	} else {
		log.Info().Msg("postcli closed")
	}

	if err = newCache.Close(ctxClose); err != nil {
		log.Error().Err(err).Msg("cache close")
		if err == context.DeadlineExceeded {
			return
		}
	} else {
		log.Info().Msg("cache closed")
	}
}
