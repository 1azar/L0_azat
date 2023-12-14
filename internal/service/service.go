package service

import (
	"L0_azat/internal/config"
	"L0_azat/internal/service/handlers/record"
	"L0_azat/internal/storage/postgres"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/nats-io/nats.go"
	"log/slog"
	"os"
)

type Service struct {
	conn  *nats.Conn
	store *postgres.Storage
	log   *slog.Logger
	cache *lru.Cache[string, any]
}

func (s *Service) Terminate() {
	s.conn.Close()
}

type credentials struct {
	clusterID string
	clientID  string
	subject   string
}

//type StorageModel interface {
//	CleanCacheReplicant() error
//}

func New(cfg *config.Config, storage *postgres.Storage, logger *slog.Logger) (*Service, error) {
	const fn = "service.service.New"
	log := logger.With(slog.String("fn", fn))

	// set nats connection
	creds := fetchCredentials(cfg)
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		return nil, err
	}

	// setting cache
	cache, err := lru.New[string, any](cfg.ServiceCfg.CacheSize)
	if err != nil {
		log.Error("cache initializing failed", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		return nil, err
	}

	// create service instance
	var service Service = Service{
		conn:  nc,
		store: storage,
		cache: cache,
		log:   logger,
	}

	// restoring cache
	if err := service.store.FillCache(service.cache); err != nil {
		log.Error("cache restoring error", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		return nil, err
	}

	// setting service handlers
	_, err = service.conn.Subscribe(creds.subject, record.New(logger, service.store, service.cache))
	if err != nil {
		log.Error("could not accomplish subscription to subject", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		return nil, err
	}

	return &service, nil
}

func fetchCredentials(cfg *config.Config) *credentials {
	// is fields emptiness check required? NO
	return &credentials{
		clusterID: os.Getenv(cfg.NatsCfg.ClusterIdEnv),
		clientID:  os.Getenv(cfg.NatsCfg.ClientIdENv),
		subject:   os.Getenv(cfg.NatsCfg.SubjectEnv),
	}
}

//func (s *Service) ReplicateCache() error {
//	const fn = "service.service.New"
//	log := s.log.With(slog.String("fn", fn))
//
//	//if err := s.store.CleanCacheReplicant(); err != nil {
//	//	log.Error("cant prune cache table in db", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
//	//	return err
//	//}
//
//	q := `
//	INSERT INTO %s (key, value)
//	VALUES ($1, $2)
//	`
//
//	query := fmt.Sprintf("INSERT INTO %s (key, value) VALUES ($1, $2)", storage.CACHE_REPLICA_TABLE)
//
//	for key, value := range s.cache {
//
//	}
//
//}
