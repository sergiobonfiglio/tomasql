package main

import (
	"testing"

	"github.com/sergiobonfiglio/tomasql/cmd/table-def-gen/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIgnoreUnknownTypes(t *testing.T) {
	// Test with ignore-unknown-types set to true
	t.Run("WithIgnoreUnknownTypes", func(t *testing.T) {
		container, err := setup.SetupTestContainer(nil, "unsupported_col_type_schema.sql", "postgres:latest")
		require.NoError(t, err)
		defer container.Close()

		tableDefData, err := getTableDefinitionFromTestDB(container, "testpkg", true)
		require.NoError(t, err)
		require.NotNil(t, tableDefData)
		require.Len(t, tableDefData.Tables, 1)

		table := tableDefData.Tables[0]
		assert.Equal(t, "test_table", table.SqlName)

		// Should have 3 columns (id, name, created_at) - metadata (json) should be skipped
		assert.Len(t, table.Columns, 3)

		columnNames := make([]string, len(table.Columns))
		for i, col := range table.Columns {
			columnNames[i] = col.SqlName
		}
		assert.Contains(t, columnNames, "id")
		assert.Contains(t, columnNames, "name")
		assert.Contains(t, columnNames, "created_at")
		assert.NotContains(t, columnNames, "metadata", "metadata column with unsupported JSON type should be skipped")
	})

	// Test with ignore-unknown-types set to false (should fail)
	t.Run("WithoutIgnoreUnknownTypes", func(t *testing.T) {
		container, err := setup.SetupTestContainer(nil, "unsupported_col_type_schema.sql", "postgres:latest")
		require.NoError(t, err)
		defer container.Close()

		tableDefData, err := getTableDefinitionFromTestDB(container, "testpkg", false)
		assert.Error(t, err, "Should fail when encountering unknown type with ignore-unknown-types=false")
		assert.Nil(t, tableDefData)
		assert.Contains(t, err.Error(), "unknown type")
	})
}
