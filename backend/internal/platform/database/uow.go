package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrTransactionRequired = errors.New("database transaction is required")

type contextKey struct{}

func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, contextKey{}, db)
}

func DBFromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if db, ok := ctx.Value(contextKey{}).(*gorm.DB); ok && db != nil {
		return db
	}
	return fallback
}

func DBFromTransactionContext(ctx context.Context) (*gorm.DB, error) {
	db, ok := ctx.Value(contextKey{}).(*gorm.DB)
	if !ok || db == nil {
		return nil, ErrTransactionRequired
	}
	return db, nil
}

func WithReadSnapshot(ctx context.Context, fallback *gorm.DB, fn func(context.Context) error) error {
	if fn == nil {
		return errors.New("read snapshot callback is nil")
	}
	if existing, ok := ctx.Value(contextKey{}).(*gorm.DB); ok && existing != nil {
		return fn(ctx)
	}
	if fallback == nil {
		return errors.New("read snapshot database is nil")
	}
	return fallback.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(WithDB(ctx, tx))
	}, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
}

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(context.Context) error) error {
	if u == nil || u.db == nil {
		return errors.New("unit of work database is nil")
	}
	if fn == nil {
		return fmt.Errorf("unit of work callback is nil")
	}
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(WithDB(ctx, tx))
	})
}
