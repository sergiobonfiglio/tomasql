package pgres

import (
	"testing"

	"github.com/sergiobonfiglio/tomasql"
)

func TestPostgresDialectName(t *testing.T) {
	dialect := &PostgresDialect{}
	expected := "postgres"
	got := dialect.Name()

	if got != expected {
		t.Errorf("Name() = %q, want %q", got, expected)
	}
}

func TestPostgresDialectPlaceholder(t *testing.T) {
	tests := []struct {
		position int
		want     string
	}{
		{1, "$1"},
		{2, "$2"},
		{10, "$10"},
		{100, "$100"},
		{0, "$0"},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := dialect.Placeholder(tt.position)
			if got != tt.want {
				t.Errorf("Placeholder(%d) = %q, want %q", tt.position, got, tt.want)
			}
		})
	}
}

func TestPostgresDialectImplementsDialectInterface(t *testing.T) {
	var _ tomasql.Dialect = (*PostgresDialect)(nil)
}

func TestGetDialect(t *testing.T) {
	dialect := GetDialect()

	if dialect == nil {
		t.Error("GetDialect() returned nil")
	}

	if dialect.Name() != "postgres" {
		t.Errorf("GetDialect().Name() = %q, want %q", dialect.Name(), "postgres")
	}
}

func TestSetDialect(t *testing.T) {
	// Save the current dialect to restore it later
	originalDialect := tomasql.GetDialect()
	defer tomasql.SetDialect(originalDialect)

	// Set the postgres dialect
	SetDialect()

	// Verify it was set correctly
	currentDialect := tomasql.GetDialect()
	if currentDialect.Name() != "postgres" {
		t.Errorf("After SetDialect(), dialect name = %q, want %q", currentDialect.Name(), "postgres")
	}

	// Test placeholder generation
	if currentDialect.Placeholder(1) != "$1" {
		t.Errorf("After SetDialect(), placeholder = %q, want %q", currentDialect.Placeholder(1), "$1")
	}
}

func TestPostgresDialectMultiplePlaceholders(t *testing.T) {
	dialect := &PostgresDialect{}

	placeholders := []string{
		dialect.Placeholder(1),
		dialect.Placeholder(2),
		dialect.Placeholder(3),
	}

	expected := []string{"$1", "$2", "$3"}

	for i, got := range placeholders {
		if got != expected[i] {
			t.Errorf("Placeholder(%d) = %q, want %q", i+1, got, expected[i])
		}
	}
}
