package gowhere

// Config defines the config for planner
type Config struct {
	// The separator between field and operation. Default to "__" which requires the condition format as: field__operator
	Separator string
	// Collections of methods to build correct SQL clause for specific dialect. Only support MySQL and PostgreSQL (default) for the moment
	Dialect Dialect
	// Whether to report error or silently skip anomalies in the conditions schema. Default to false
	Strict bool
	// Table name to add before the columns in SQL clause, i.e: table_name.column_name. Default to empty which will keep the column unchanged
	Table string
	// The map of column aliases to be replaced when build the SQL clause. Use cases:
	// Example: {"name": "foo.name", "bname": "bar.name"}
	ColumnAliases map[string]string
	// Custom conditions allow full access on the condition generating
	CustomConditions map[string]CustomConditionFn

	// sort the conditions in map by field for testing purposes only
	sort bool
}

var (
	// DefaultConfig is the default configuration of the planner
	DefaultConfig = Config{
		Separator:        "__",
		Dialect:          DialectPostgreSQL,
		ColumnAliases:    make(map[string]string),
		CustomConditions: make(map[string]CustomConditionFn),
		// Strict: false,
		// Table: "",
	}
)

// Modes to set configurations
const (
	OverwriteMode = 'o'
	AppendMode    = 'a'
	WriteMode     = 'w'
)
