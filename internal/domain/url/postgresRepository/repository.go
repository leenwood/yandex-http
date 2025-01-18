package postgresRepository

import (
	"context"
	"database/sql"
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
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.Username,
		config.Password,
		config.Hostname,
		config.Port,
		config.Database,
		config.SSLMode,
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
	err = r.db.QueryRow(r.ctx, query, args...).
		Scan(&model.Id, &model.OriginalUrl, &model.ClickCount, &model.CreatedDate)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Возвращаем nil вместо ошибки
		}
		return nil, err
	}

	model := &url.Url{}
	err = r.db.QueryRow(r.ctx, query, args...).
		Scan(&model.Id, &model.OriginalUrl, &model.ClickCount, &model.CreatedDate)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// Save метод для сохранения ссылки, второй параметр овтечает за ИД ссылки
// Если передать конкретный ИД и он будет занят, вернется ошибка
func (r *Repository) Save(originalUrl, shortUuid string) (*url.Url, error) {
	var err error
	if shortUuid == "" {
		shortUuid, err = r.GenerateUuid()
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

	if err != nil {
		return nil, err
	}
	query, args, err := r.sq.
		Insert("urls").
		Columns("id", "original_url", "click_count", "created_date").
		Values(shortUuid, originalUrl, 0, createdDate).
		ToSql()
	var model *url.Url
	_, err = r.db.Exec(r.ctx, query, args)
	if err != nil {
		return nil, err
	}

	// Возвращаем сущность с сгенерированным id
	// Предпологаем что если ошибки не было то ИД не занят и сохранился с тем что создался
	model.Id = shortUuid
	model.OriginalUrl = originalUrl
	model.ClickCount = 0
	model.CreatedDate = createdDate
	return model, nil
}

// GenerateUuid Метод для генерации UUID которые не заняты
// Генерируем на стороне кода чтобы была возможность делать кастомные UUID
func (r *Repository) GenerateUuid() (string, error) {
	var shortUUID string
	var exists bool
	var err error
	for {
		select {
		case <-r.ctx.Done():
			break
		default:
			shortUUID = uuid.New().String()[:5]
			exists, err = r.IsIdExists(shortUUID)
			if err != nil {
				return "", err
			}
			if !exists {
				break
			}
		}
	}
}

// IsIdExists Функция проверки уникальности индефикатора
// true - UUID свободен
// false - UUID занят
func (r *Repository) IsIdExists(id string) (bool, error) {
	query, _, err := r.sq.
		Select("EXISTS (SELECT 1 FROM urls WHERE id = $1)").
		ToSql()
	if err != nil {
		return false, err
	}
	var exists bool
	err = r.db.QueryRow(r.ctx, query, id).Scan(&exists)
	if err != nil && err != pgx.ErrNoRows {
		return false, err
	}
	return exists, nil
}
