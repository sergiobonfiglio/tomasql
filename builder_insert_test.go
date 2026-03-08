package tomasql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilderInsert(t *testing.T) {
	t.Run("single row insert", func(t *testing.T) {
		sql, params := InsertInto(Account).
			Columns(Account.Uuid, Account.Type).
			ValuesParam([]any{"test-uuid-123", "user"}).
			SQL()

		require.Equal(t, "INSERT INTO account (uuid, type) VALUES ("+GetDialect().Placeholder(1)+", "+GetDialect().Placeholder(2)+")", sql)
		require.Equal(t, []any{"test-uuid-123", "user"}, params)
	})

	t.Run("batch insert", func(t *testing.T) {
		sql, params := InsertInto(Account).
			Columns(Account.Uuid, Account.Type).
			ValuesParam([]any{"test-uuid-123", "user"}).
			ValuesParam([]any{"test-uuid-456", "admin"}).
			SQL()

		require.Equal(t, "INSERT INTO account (uuid, type) VALUES ("+GetDialect().Placeholder(1)+", "+GetDialect().Placeholder(2)+"), ("+GetDialect().Placeholder(3)+", "+GetDialect().Placeholder(4)+")", sql)
		require.Equal(t, []any{"test-uuid-123", "user", "test-uuid-456", "admin"}, params)
	})

	t.Run("single row with all columns", func(t *testing.T) {
		sql, params := InsertInto(Account).
			Columns(Account.Id, Account.Uuid, Account.Type, Account.CreatedTs).
			ValuesParam([]any{int64(1), "test-uuid-123", "user", 1234567890}).
			SQL()

		require.Equal(t, "INSERT INTO account (id, uuid, type, created_ts) VALUES ("+GetDialect().Placeholder(1)+", "+GetDialect().Placeholder(2)+", "+GetDialect().Placeholder(3)+", "+GetDialect().Placeholder(4)+")", sql)
		require.Equal(t, []any{int64(1), "test-uuid-123", "user", 1234567890}, params)
	})

	t.Run("batch with three rows", func(t *testing.T) {
		sql, params := InsertInto(Account).
			Columns(Account.Uuid, Account.Type).
			ValuesParam([]any{"uuid-1", "type-1"}).
			ValuesParam([]any{"uuid-2", "type-2"}).
			ValuesParam([]any{"uuid-3", "type-3"}).
			SQL()

		require.Equal(t, "INSERT INTO account (uuid, type) VALUES ("+GetDialect().Placeholder(1)+", "+GetDialect().Placeholder(2)+"), ("+GetDialect().Placeholder(3)+", "+GetDialect().Placeholder(4)+"), ("+GetDialect().Placeholder(5)+", "+GetDialect().Placeholder(6)+")", sql)
		require.Equal(t, []any{"uuid-1", "type-1", "uuid-2", "type-2", "uuid-3", "type-3"}, params)
	})
}
