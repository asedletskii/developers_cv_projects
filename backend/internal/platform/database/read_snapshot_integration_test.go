//go:build integration

package database

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/testutil/pgtest"
)

func TestWithReadSnapshotUsesConsistentReadOnlyRepeatableReadTransaction(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test")
	}

	ctx, dsn := pgtest.NewDatabase(t)

	db, err := NewGORM(dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, Close(db)) })
	require.NoError(t, db.Exec("CREATE SCHEMA IF NOT EXISTS test_probe").Error)
	require.NoError(t, db.Exec("CREATE TABLE test_probe.snapshot_probe (id integer PRIMARY KEY, value integer NOT NULL)").Error)
	require.NoError(t, db.Exec("INSERT INTO test_probe.snapshot_probe (id, value) VALUES (1, 10)").Error)

	writer, err := pgx.Connect(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, writer.Close(ctx)) })

	var before int
	var after int
	var isolation string
	var readOnly string
	err = WithReadSnapshot(ctx, db, func(snapshotCtx context.Context) error {
		tx, err := DBFromTransactionContext(snapshotCtx)
		if err != nil {
			return err
		}
		if err := tx.Raw("SHOW transaction_isolation").Scan(&isolation).Error; err != nil {
			return err
		}
		if err := tx.Raw("SHOW transaction_read_only").Scan(&readOnly).Error; err != nil {
			return err
		}
		if err := tx.Raw("SELECT value FROM test_probe.snapshot_probe WHERE id = 1").Scan(&before).Error; err != nil {
			return err
		}
		if _, err := writer.Exec(snapshotCtx, "UPDATE test_probe.snapshot_probe SET value = 20 WHERE id = 1"); err != nil {
			return err
		}
		return tx.Raw("SELECT value FROM test_probe.snapshot_probe WHERE id = 1").Scan(&after).Error
	})
	require.NoError(t, err)
	require.Equal(t, "repeatable read", isolation)
	require.Equal(t, "on", readOnly)
	require.Equal(t, 10, before)
	require.Equal(t, 10, after, "the transaction must not observe a commit that happened after its snapshot")

	var committed int
	require.NoError(t, db.Raw("SELECT value FROM test_probe.snapshot_probe WHERE id = 1").Scan(&committed).Error)
	require.Equal(t, 20, committed)
}
