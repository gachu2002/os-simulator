package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"os-simulator-plan/internal/db/sqlc"
	"os-simulator-plan/internal/lessons"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const lessonProgressMetadataKey = "lessons.progress.v1"

type metadataProgressPersistence struct {
	queries *sqlc.Queries
}

func NewLessonEngineWithPersistence(pool *pgxpool.Pool) *lessons.Engine {
	if pool == nil {
		return lessons.NewEngine()
	}
	persistence := &metadataProgressPersistence{queries: sqlc.New(pool)}
	return lessons.NewEngineWithCatalogAndPersistence(lessons.DefaultCatalog(), persistence)
}

func (p *metadataProgressPersistence) Load() (map[string]lessons.StageProgress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	row, err := p.queries.GetMetadata(ctx, lessonProgressMetadataKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return map[string]lessons.StageProgress{}, nil
		}
		return nil, err
	}

	out := map[string]lessons.StageProgress{}
	if err := json.Unmarshal([]byte(row.Value), &out); err != nil {
		return nil, fmt.Errorf("decode lesson progress: %w", err)
	}
	return out, nil
}

func (p *metadataProgressPersistence) Save(stages map[string]lessons.StageProgress) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	payload, err := json.Marshal(stages)
	if err != nil {
		return fmt.Errorf("encode lesson progress: %w", err)
	}

	return p.queries.UpsertMetadata(ctx, sqlc.UpsertMetadataParams{
		Key:   lessonProgressMetadataKey,
		Value: string(payload),
	})
}
