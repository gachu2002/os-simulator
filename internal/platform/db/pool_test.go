package db

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestNewPoolFromEnvSkipsWhenDatabaseURLUnset(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	pool, err := NewPoolFromEnv(context.Background(), zap.NewNop())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if pool != nil {
		t.Fatalf("expected nil pool when DATABASE_URL is unset")
	}
}

func TestNewPoolFromEnvSkipsWhenDatabaseURLWhitespace(t *testing.T) {
	t.Setenv("DATABASE_URL", "   ")

	pool, err := NewPoolFromEnv(context.Background(), zap.NewNop())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if pool != nil {
		t.Fatalf("expected nil pool when DATABASE_URL is blank")
	}
}
