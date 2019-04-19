package gowhere

// WithConfig returns an empty plan using the given configs. Zero value will be replaced by default config.
func WithConfig(conf Config) *Plan {
	if conf.Separator == "" {
		conf.Separator = DefaultConfig.Separator
	}
	if conf.Dialect == nil {
		conf.Dialect = DefaultConfig.Dialect
	}
	if conf.ColumnAliases == nil {
		conf.ColumnAliases = DefaultConfig.ColumnAliases
	}

	return &Plan{config: &conf, conditions: &andConditions{naked: true}}
}

// Where is shortcut to create new plan with default configurations
func Where(cond interface{}, vars ...interface{}) *Plan {
	return WithConfig(DefaultConfig).Where(cond, vars...)
}

// WhereMySQL returns a plan with given conditions for MySQL, using default configurations
func WhereMySQL(cond interface{}, vars ...interface{}) *Plan {
	return WithConfig(Config{Dialect: DialectMySQL}).Where(cond, vars...)
}

// WherePostgreSQL returns a plan with given conditions for PostgreSQL, using default configurations
func WherePostgreSQL(cond interface{}, vars ...interface{}) *Plan {
	return WithConfig(Config{Dialect: DialectPostgreSQL}).Where(cond, vars...)
}
