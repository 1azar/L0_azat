package postgres

import (
	"L0_azat/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

type Storage struct {
	conn *pgx.Conn
}

type credentials struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

func New(cfg *config.Config) (*Storage, error) {
	creds := fetchCredentials(cfg)

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		creds.host, creds.port, creds.user, creds.password, creds.dbname)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &Storage{conn: conn}, nil
}

func (s *Storage) Close() {
	s.conn.Close(context.Background())
}

func fetchCredentials(cfg *config.Config) *credentials {
	var credentials *credentials = &credentials{
		host:     os.Getenv(cfg.DbCredentials.AddressEnv),
		port:     os.Getenv(cfg.DbCredentials.PortEnv),
		user:     os.Getenv(cfg.DbCredentials.UsernameEnv),
		password: os.Getenv(cfg.DbCredentials.PasswordEnv),
		dbname:   os.Getenv(cfg.DbCredentials.DbNameEnv),
	}
	// is fields emptiness check required? NO
	return credentials
}
