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
		conf.ColumnAliases = make(map[string]string)
		for key, val := range DefaultConfig.ColumnAliases {
			conf.ColumnAliases[key] = val
		}
	}
	if conf.CustomConditions == nil {
		conf.CustomConditions = make(map[string]CustomConditionFn)
		for key, val := range DefaultConfig.CustomConditions {
			conf.CustomConditions[key] = val
		}
	}

	return &Plan{config: &conf, conditions: &andConditions{naked: true}}
}

// Where is shortcut to create new plan with default configurations
func Where(cond interface{}, vars ...interface{}) *Plan {
	return WithConfig(Config{}).Where(cond, vars...)
}
