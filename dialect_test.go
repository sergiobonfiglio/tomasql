package tomasql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestStandardDialect_Name tests the Name method
func TestStandardDialect_Name(t *testing.T) {
	dialect := &standardDialect{}
	require.Equal(t, "standard", dialect.Name())
}

// TestStandardDialect_Placeholder tests placeholder generation
func TestStandardDialect_Placeholder(t *testing.T) {
	tests := []struct {
		position int
		expected string
	}{
		{1, "?"},
		{2, "?"},
		{100, "?"},
	}

	dialect := &standardDialect{}
	for _, tt := range tests {
		t.Run("Position "+string(rune(tt.position)), func(t *testing.T) {
			require.Equal(t, tt.expected, dialect.Placeholder(tt.position))
		})
	}
}

// TestGetDialect tests that GetDialect returns the current dialect
func TestGetDialect(t *testing.T) {
	// Reset to default
	SetDialect(DefaultDialect)

	d := GetDialect()
	require.NotNil(t, d)
	require.Equal(t, "standard", d.Name())
	require.Equal(t, "?", d.Placeholder(1))
}

// TestSetDialect_ChangeDialect tests that SetDialect changes the dialect
func TestSetDialect_ChangeDialect(t *testing.T) {
	// Save original dialect
	originalDialect := GetDialect()
	defer SetDialect(originalDialect)

	// Create a custom dialect
	customDialect := &customTestDialect{name: "custom"}

	// Change dialect
	SetDialect(customDialect)

	// Verify it changed
	d := GetDialect()
	require.Equal(t, "custom", d.Name())
}

// TestSetDialect_MultipleChanges tests multiple dialect changes
func TestSetDialect_MultipleChanges(t *testing.T) {
	// Save original dialect
	originalDialect := GetDialect()
	defer SetDialect(originalDialect)

	// Change to first custom dialect
	dialect1 := &customTestDialect{name: "dialect1"}
	SetDialect(dialect1)
	require.Equal(t, "dialect1", GetDialect().Name())

	// Change to second custom dialect
	dialect2 := &customTestDialect{name: "dialect2"}
	SetDialect(dialect2)
	require.Equal(t, "dialect2", GetDialect().Name())

	// Change back to first
	SetDialect(dialect1)
	require.Equal(t, "dialect1", GetDialect().Name())
}

// TestDefaultDialect tests the default dialect is correct
func TestDefaultDialect(t *testing.T) {
	require.NotNil(t, DefaultDialect)
	require.Equal(t, "standard", DefaultDialect.Name())
	require.Equal(t, "?", DefaultDialect.Placeholder(1))
}

// TestDialect_InterfaceCompliance tests that standardDialect implements Dialect
func TestDialect_InterfaceCompliance(t *testing.T) {
	var _ Dialect = (*standardDialect)(nil)
}

// TestDialect_Placeholder_IgnoresPosition tests that standard dialect ignores position
func TestDialect_Placeholder_IgnoresPosition(t *testing.T) {
	dialect := &standardDialect{}

	for i := 1; i <= 10; i++ {
		require.Equal(t, "?", dialect.Placeholder(i))
	}
}

// customTestDialect is a test dialect for testing SetDialect
type customTestDialect struct {
	name string
}

var _ Dialect = (*customTestDialect)(nil)

func (d *customTestDialect) Name() string {
	return d.name
}

func (d *customTestDialect) Placeholder(_ int) string {
	return "custom"
}
