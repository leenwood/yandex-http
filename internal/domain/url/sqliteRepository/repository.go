package sqliteRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"leenwood/yandex-http/config"
	"leenwood/yandex-http/internal/domain/url"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db  *sql.DB
	sq  sq.StatementBuilderType
	ctx context.Context
}

func NewRepository(ctx context.Context, _ config.DatabaseConfig) (*Repository, error) {
	dsn := "database.sqlite"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return &Repository{
		db:  db,
		sq:  sq.StatementBuilder.PlaceholderFormat(sq.Question),
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
	row := r.db.QueryRowContext(r.ctx, query, args...)
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
	row := r.db.QueryRowContext(r.ctx, query, args...)
	err = row.Scan(&model.Id, &model.OriginalUrl, &model.ClickCount, &model.CreatedDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Возвращаем nil вместо ошибки
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

	_, err = r.db.ExecContext(r.ctx, query, args...)
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
	query, _, err := r.sq.
		Select("EXISTS (SELECT 1 FROM urls WHERE id = ?)").
		ToSql()
	if err != nil {
		return false, err
	}

	var exists bool
	row := r.db.QueryRowContext(r.ctx, query, id)
	err = row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}
