package postgres

import (
	"L0_azat/internal/config"
	"L0_azat/internal/domain"
	"L0_azat/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
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

func (s *Storage) GetOrder(orderUid string) (*domain.Message, error) {
	const fn = "storage.postgres.GetOrder"

	// todo: try use sql injection when this case:
	//q := fmt.Sprintf("SELECT * FROM %s WHERE order_uid = %s LIMIT 1", storage.ORDERS_TABLE, orderUid)

	// get order without items
	q := fmt.Sprintf("SELECT * FROM %s WHERE order_uid = $1 LIMIT 1", storage.ORDERS_TABLE)

	row := s.conn.QueryRow(context.Background(), q, orderUid)

	var msg domain.Message
	err := row.Scan(
		&msg.OrderUid,
		&msg.TrackNumber,
		&msg.Entry,
		&msg.DeliveryInfo,
		&msg.PaymentInfo,
		&msg.Locale,
		&msg.InternalSignature,
		&msg.CustomerId,
		&msg.DeliveryService,
		&msg.Shardkey,
		&msg.SmId,
		&msg.DateCreated,
		&msg.OofShard,
	)
	if errors.Is(pgx.ErrNoRows, err) {
		return &domain.Message{}, storage.ErrOrderNotFound
	}
	if err != nil {
		return &domain.Message{}, fmt.Errorf("failed to get order. %s: %w", fn, err)
	}

	// get items for order
	q = fmt.Sprintf("SELECT * FROM %s WHERE order_uid = $1", storage.ORDER_ITEMS_TABLE)
	rows, err := s.conn.Query(context.Background(), q, orderUid)
	if err != nil {
		return &domain.Message{}, fmt.Errorf("failed to get order items. %s: %w", fn, err)
	}
	defer rows.Close()

	var items []domain.Item = make([]domain.Item, 0, 10)
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(
			&item.ChrtId,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmId,
			&item.Brand,
			&item.Status,
		); err != nil {
			return &domain.Message{}, fmt.Errorf("failed to get order items. %s: %w", fn, err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return &domain.Message{}, fmt.Errorf("%s: %w", fn, err)
	}

	msg.Items = items

	return &msg, nil
}

func (s *Storage) FillCache(c *lru.Cache[string, any]) error {
	const fn = "storage.postgres.FillCache"

	q := fmt.Sprintf("SELECT * FROM %s LIMIT %d", storage.ORDERS_TABLE, c.Len())

	rows, err := s.conn.Query(context.Background(), q)
	if err != nil {
		return fmt.Errorf("failed to get rows in orders table. %s: %w", fn, err)
	}
	defer rows.Close()
	for rows.Next() {
		var value domain.Message
		if err := rows.Scan(
			&value.OrderUid,
			&value.TrackNumber,
			&value.Entry,
			&value.DeliveryInfo,
			&value.PaymentInfo,
			&value.Locale,
			&value.InternalSignature,
			&value.CustomerId,
			&value.DeliveryService,
			&value.Shardkey,
			&value.SmId,
			&value.DateCreated,
			&value.OofShard,
		); err != nil {
			return fmt.Errorf("failed to get row. %s: %w", fn, err)
		}
		c.Add(value.OrderUid, value)
	}
	if rows.Err() != nil {
		return fmt.Errorf("rows error. %s: %w", fn, err)
	}

	return nil
}

//func (s *Storage) CleanCacheReplicant() error {
//	const fn = "storage.postgres.CleanCacheReplicant"
//
//	q := fmt.Sprintf("DELETE FROM %s", storage.CACHE_REPLICA_TABLE)
//	_, err := s.conn.Exec(context.Background(), q)
//	if err != nil {
//		return fmt.Errorf("failed to prune cache replica table (%s). %s: %w", storage.CACHE_REPLICA_TABLE, fn, err)
//	}
//	return nil
//}

func (s *Storage) SaveMsg(msg domain.Message) error {
	const fn = "storage.postgres.SaveMsg"

	deliveryInfo, err := json.Marshal(msg.DeliveryInfo)
	if err != nil {
		return fmt.Errorf("marshal delivery info failed. %s: %w", fn, err)
	}

	paymentInfo, err := json.Marshal(msg.PaymentInfo)
	if err != nil {
		return fmt.Errorf("marshal payment info failed. %s: %w", fn, err)
	}

	// ORDERS table
	query := fmt.Sprintf(`INSERT INTO %s
							(order_uid,
							 track_number,
							 entry,
							 delivery_info,
							 payment_info,
							 locale,
							 internal_signature,
							 customer_id,
							 delivery_service,
							 shardkey,
							 sm_id,
							 date_created,
							 oof_shard
							) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ,$11 ,$12 ,$13)`, storage.ORDERS_TABLE)
	_, err = s.conn.Exec(context.Background(), query,
		msg.OrderUid,
		msg.TrackNumber,
		msg.Entry,
		deliveryInfo,
		paymentInfo,
		msg.Locale,
		msg.InternalSignature,
		msg.CustomerId,
		msg.DeliveryService,
		msg.Shardkey,
		msg.SmId,
		msg.DateCreated,
		msg.OofShard,
	)
	if err != nil {
		return fmt.Errorf("could not save order info. %s: %w", fn, err)
	}

	// ORDER_ITEMS table
	for idx, item := range msg.Items {
		q := fmt.Sprintf(`INSERT INTO %s (
												order_uid,
												chrt_id,
												price,
												rid,
												name,
												sale,
												size,
												total_price,
												nm_id,
												brand,
												status
                ) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, storage.ORDER_ITEMS_TABLE)
		_, err = s.conn.Exec(context.Background(), q, msg.OrderUid,
			item.ChrtId,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmId,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return fmt.Errorf("could not save order items info. %s: %w. %d item of %d: %v", fn, err, idx, len(msg.Items), item)
		}
	}

	return nil
}

func fetchCredentials(cfg *config.Config) *credentials {
	// is fields emptiness check required? NO
	return &credentials{
		host:     os.Getenv(cfg.DbCredentials.AddressEnv),
		port:     os.Getenv(cfg.DbCredentials.PortEnv),
		user:     os.Getenv(cfg.DbCredentials.UsernameEnv),
		password: os.Getenv(cfg.DbCredentials.PasswordEnv),
		dbname:   os.Getenv(cfg.DbCredentials.DbNameEnv),
	}
}
