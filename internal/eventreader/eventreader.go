package eventreader

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/v1tbrah/feed-service/config"
)

type Reader struct {
	liveCtx context.Context

	readerPostCreated *kafka.Reader
	readerPostDeleted *kafka.Reader

	readerFriendAdded   *kafka.Reader
	readerFriendRemoved *kafka.Reader

	cache Cache

	relationCli RelationCli

	postCli PostCli

	cfg config.Kafka
}

func Init(ctx context.Context, cache Cache, relationCli RelationCli, postCli PostCli, cfg config.Kafka) (*Reader, error) {
	if !cfg.Enable {
		return nil, nil
	}

	brokerAddr := net.JoinHostPort(cfg.Host, cfg.Port)

	reader := &Reader{
		liveCtx: ctx,

		readerPostCreated: kafka.NewReader(kafka.ReaderConfig{Brokers: []string{brokerAddr}, Topic: cfg.TopicPostCreated}),
		readerPostDeleted: kafka.NewReader(kafka.ReaderConfig{Brokers: []string{brokerAddr}, Topic: cfg.TopicPostDeleted}),

		readerFriendAdded:   kafka.NewReader(kafka.ReaderConfig{Brokers: []string{brokerAddr}, Topic: cfg.TopicFriendAdded}),
		readerFriendRemoved: kafka.NewReader(kafka.ReaderConfig{Brokers: []string{brokerAddr}, Topic: cfg.TopicFriendRemoved}),

		cache: cache,

		relationCli: relationCli,

		postCli: postCli,

		cfg: cfg,
	}

	go func() {
		log.Info().Msg("start read topic post created")
		reader.startProcessEventPostCreated()
	}()

	go func() {
		log.Info().Msg("start read topic post deleted")
		reader.startProcessEventPostDeleted()
	}()

	go func() {
		log.Info().Msg("start read topic friend added")
		reader.startProcessEventFriendAdded()
	}()

	go func() {
		log.Info().Msg("start read topic friend removed")
		reader.startProcessEventFriendRemoved()
	}()

	return reader, nil
}

func (mr *Reader) Close(ctx context.Context) (err error) {
	if mr == nil {
		return nil
	}

	closed := make(chan struct{})

	go func() {
		if err = mr.readerPostCreated.Close(); err != nil {
			err = errors.Wrap(err, "readerPostCreated.Close")
			closed <- struct{}{}
			return
		}

		if err = mr.readerPostDeleted.Close(); err != nil {
			err = errors.Wrap(err, "readerPostDeleted.Close")
			closed <- struct{}{}
			return
		}

		closed <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-closed:
		return err
	}
}
