package postgresRepository

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"leenwood/yandex-http/config"
	"leenwood/yandex-http/internal/domain/url"
	"time"
)

type Repository struct {
	db  *pgxpool.Pool
	sq  sq.StatementBuilderType
	ctx context.Context
}

func NewRepository(ctx context.Context, config config.DatabaseConfig) (*Repository, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.Username,
		config.Password,
		config.Hostname,
		config.Port,
		config.Database,
	)
	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Repository{
		db:  dbpool,
		sq:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		ctx: ctx,
	}, nil
}

func (r *Repository) FindById(id string) (*url.Url, error) {
	query, args, err := r.sq.
		Select("id", "original_url", "click_count", "created_date").
		From("urls").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	model := &url.Url{}
	row := r.db.QueryRow(r.ctx, query, args...)
	err = row.Scan(&model.Id, &model.OriginalUrl, &model.ClickCount, &model.CreatedDate)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *Repository) FindByUrl(originalUrl string) (*url.Url, error) {
	query, args, err := r.sq.
		Select("id", "original_url", "click_count", "created_date").
		From("urls").
		Where(sq.Eq{"original_url": originalUrl}).
		ToSql()
	if err != nil {
		return nil, err
	}

	model := &url.Url{}
	row := r.db.QueryRow(r.ctx, query, args...)
	err = row.Scan(&model.Id, &model.OriginalUrl, &model.ClickCount, &model.CreatedDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return model, nil
}

func (r *Repository) Save(originalUrl, shortUuid string) (*url.Url, error) {
	var err error
	if shortUuid == "" {
		shortUuid, err = r.GenerateUuid()
		if err != nil {
			return nil, err
		}
	} else {
		isExists, err := r.IsIdExists(shortUuid)
		if err != nil {
			return nil, err
		}
		if isExists {
			return nil, fmt.Errorf("short uuid already exists")
		}
	}

	createdDate := time.Now()
	query, args, err := r.sq.
		Insert("urls").
		Columns("id", "original_url", "click_count", "created_date").
		Values(shortUuid, originalUrl, 0, createdDate).
		ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(r.ctx, query, args...)
	if err != nil {
		return nil, err
	}

	model := &url.Url{
		Id:          shortUuid,
		OriginalUrl: originalUrl,
		ClickCount:  0,
		CreatedDate: createdDate,
	}
	return model, nil
}

func (r *Repository) GenerateUuid() (string, error) {
	for {
		select {
		case <-r.ctx.Done():
			return "", r.ctx.Err()
		default:
			shortUUID := uuid.New().String()[:5]
			exists, err := r.IsIdExists(shortUUID)
			if err != nil {
				return "", err
			}
			if !exists {
				return shortUUID, nil
			}
		}
	}
}

func (r *Repository) IsIdExists(id string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM urls WHERE id = $1)"

	var exists bool
	err := r.db.QueryRow(r.ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) FindAll(page, limit int) ([]*url.Url, error) {
	offset := (page - 1) * limit

	query, args, err := r.sq.
		Select("id", "original_url", "click_count", "created_date").
		From("urls").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(r.ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []*url.Url
	for rows.Next() {
		var u url.Url
		if err := rows.Scan(&u.Id, &u.OriginalUrl, &u.ClickCount, &u.CreatedDate); err != nil {
			return nil, err
		}
		urls = append(urls, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (r *Repository) Update(shortUrl *url.Url) (*url.Url, error) {
	if shortUrl == nil {
		return nil, errors.New("input URL cannot be nil")
	}

	query, args, err := r.sq.
		Update("urls").
		Set("original_url", shortUrl.OriginalUrl).
		Set("click_count", shortUrl.ClickCount).
		Set("created_date", shortUrl.CreatedDate).
		Where(sq.Eq{"id": shortUrl.Id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = r.db.Exec(r.ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update query: %w", err)
	}

	return shortUrl, nil
}
