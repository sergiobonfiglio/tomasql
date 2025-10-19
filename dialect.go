package tomasql

type Dialect interface {

	// Name returns the name of the dialect (e.g., "standard", "postgres")
	Name() string

	// Placeholder returns the parameter placeholder for position n (e.g., $1, ?, :1)
	Placeholder(position int) string
}

// DefaultDialect is used when no dialect is specified
var DefaultDialect Dialect = &standardDialect{}

var dialect Dialect = DefaultDialect

// SetDialect allows changing the dialect
func SetDialect(d Dialect) {
	dialect = d
}

// GetDialect returns the current dialect
func GetDialect() Dialect {
	return dialect
}

type standardDialect struct {
}

var _ Dialect = (*standardDialect)(nil)


func (d *standardDialect) Name() string {
	return "standard"
}

func (d *standardDialect) Placeholder(_ int) string {
	return "?"
}