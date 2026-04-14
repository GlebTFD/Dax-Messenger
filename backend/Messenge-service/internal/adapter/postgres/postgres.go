package postgres

import (
	"context"
	"fmt"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/dto"
	"github.com/google/uuid"

	"github.com/hashicorp/go-hclog"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	DBName   string `env:"DB_NAME"`
}

type Pool struct {
	log  hclog.Logger
	pool *pgxpool.Pool
}

func New(ctx context.Context, log hclog.Logger, cfg Config) (*Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("error to create a pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Pool{pool: pool, log: log}, nil
}

func (p *Pool) CreateMessage(ctx context.Context, msg *dto.MessageJSON) error {
	// generate uuid
	id := uuid.New()

	_, err := p.pool.Exec(ctx,
		"INSERT INTO messages (id, type, timestamp, text, reply_to, sended_id) VALUES ($1, $2, $3, $4, $5, $6)",
		id, msg.Type, msg.Timestamp, msg.Payload.Text, msg.Payload.ReplyTo, msg.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to write message to db: %w", err)
	}

	return nil
}

func (p *Pool) DeleteMessage(ctx context.Context, msgId string) (string, error) {
	var replyTo string
	err := p.pool.QueryRow(ctx,
		"DELETE FROM messages WHERE id=$1 RETURNING reply_to", msgId,
	).Scan(&replyTo)

	if err != nil {
		return "", fmt.Errorf("failed to delete message(db): %w", err)
	}

	return replyTo, nil
}

func (p *Pool) UpdateMessage(ctx context.Context, msgId string, text string) (string, error) {
	var replyTo string
	err := p.pool.QueryRow(ctx,
		"UPDATE messages SET text=$2 WHERE id=$1 RETURNING reply_to", msgId, text,
	).Scan(&replyTo)

	if err != nil {
		return "", fmt.Errorf("failed to update message(db): %w", err)
	}

	return replyTo, nil
}
