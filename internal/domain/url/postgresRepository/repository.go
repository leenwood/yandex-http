package postgresRepository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"leenwood/yandex-http/internal/domain/url"
)

type Repository struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func NewRepository(ctx context.Context, dsn string) (*Repository, error) {
	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Repository{
		db:  dbpool,
		ctx: ctx,
	}, nil
}

func (r *Repository) FindById(id string) (*url.Url, error) {
	query := `
		SELECT id, original_url, click_count, date 
		FROM urls 
		WHERE id = $1
	`

	model := &url.Url{}
	err := r.db.QueryRow(r.ctx, query, id).Scan(&model.Id, &model.OriginalUrl, &model.ClickCount, &model.Date)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return model, nil
}

func (r *Repository) FindByUrl(originalUrl string) (*url.Url, error) {
	query := `
		SELECT id, original_url, click_count, date 
		FROM urls 
		WHERE original_url = $1
	`
	model := &url.Url{}
	err := r.db.QueryRow(r.ctx, query, originalUrl).Scan(&model.Id, &model.OriginalUrl, &model.ClickCount, &model.Date)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return model, nil
}

func (r *Repository) Save(model *url.Url) (*url.Url, error) {
	// Генерируем UUID и обрезаем его до 5 символов
	var shortUUID string
	var exists bool
	var err error

	// Пытаемся генерировать уникальный UUID
	for {
		shortUUID = uuid.New().String()[:5]
		exists, err = r.isIdExists(shortUUID)
		if err != nil {
			return nil, err
		}
		if !exists {
			break
		}
	}

	// Вставляем сущность в базу данных
	query := `
		INSERT INTO urls (id, original_url, click_count, date)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var generatedId string
	err = r.db.QueryRow(r.ctx, query, shortUUID, model.OriginalUrl, model.ClickCount, model.Date).Scan(&generatedId)
	if err != nil {
		return nil, err
	}

	// Возвращаем сущность с сгенерированным id
	model.Id = generatedId
	return model, nil
}

func (r *Repository) isIdExists(id string) (bool, error) {
	query := `
		SELECT 1
		FROM urls 
		WHERE id = $1
	`
	var exists bool
	err := r.db.QueryRow(r.ctx, query, id).Scan(&exists)
	if err != nil && err != pgx.ErrNoRows {
		return false, err
	}
	return exists, nil
}
