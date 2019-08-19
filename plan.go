package gowhere

// Plan contains information to build WHERE clause
type Plan struct {
	Error error

	conditions *andConditions
	config     *Config
	built      bool
	sql        string
	vars       []interface{}
}

// Where adds more condition(s) to the current Plan, using AND operator
func (p *Plan) Where(cond interface{}, vars ...interface{}) *Plan {
	condition, err := toCondition(cond, vars, false)
	if err != nil {
		if p.config.Strict {
			p.Error = &InvalidCond{cond, vars}
		}
		return p
	}

	p.conditions.value = append(p.conditions.value, condition)
	p.built = false

	return p
}

// And is the alias of Where
func (p *Plan) And(cond interface{}, vars ...interface{}) *Plan {
	return p.Where(cond, vars...)
}

// Or wraps all current conditions and ties with the new "cond" by OR operator
func (p *Plan) Or(cond interface{}, vars ...interface{}) *Plan {
	condition, err := toCondition(cond, vars, false)
	if err != nil {
		if p.config.Strict {
			p.Error = &InvalidCond{cond, vars}
		}
		return p
	}

	p.conditions.naked = false
	p.conditions = &andConditions{naked: true, value: []interface{}{&orConditions{value: []interface{}{p.conditions, condition}}}}
	p.built = false

	return p
}

// Not works similar to Where but reverses the condition operator(s)
func (p *Plan) Not(cond interface{}, vars ...interface{}) *Plan {
	condition, err := toCondition(cond, vars, true)
	if err != nil {
		if p.config.Strict {
			p.Error = &InvalidCond{cond, vars}
		}
		return p
	}

	p.conditions.value = append(p.conditions.value, condition)
	p.built = false

	return p
}

// Build builds the SQL clause and vars with given conditions
func (p *Plan) Build() (rp *Plan) {
	defer func() {
		if err := recover(); err != nil {
			switch e := err.(type) {
			case *InvalidCond:
				// better to double check the Strict config
				if p.config.Strict {
					p.Error = e
				}
			default:
				// critical error
				p.Error = err.(error)
			}
			p.built = false
			rp = p
		}
	}()

	p.sql, p.vars = p.conditions.build(p.config)
	p.built = true

	return p
}

// SQL returns the built SQL clause
func (p *Plan) SQL() string {
	if !p.built {
		p.Build()
	}
	return p.sql
}

// Vars returns the list of vars for the built SQL clause
func (p *Plan) Vars() []interface{} {
	if !p.built {
		p.Build()
	}
	return p.vars
}

// SetTable updates the `Table` config value
func (p *Plan) SetTable(value string) *Plan {
	p.config.Table = value
	p.built = false
	return p
}

// SetColumnAliases updates the `ColumnAliases` config value
func (p *Plan) SetColumnAliases(aliases map[string]string, mode ...rune) *Plan {
	m := AppendMode
	if len(mode) > 0 && (mode[0] == OverwriteMode || mode[0] == WriteMode) {
		m = mode[0]
	}

	if m == OverwriteMode {
		p.config.ColumnAliases = aliases
	} else {
		for key, val := range aliases {
			if _, ok := p.config.ColumnAliases[key]; ok && m == AppendMode {
				continue
			}
			p.config.ColumnAliases[key] = val
		}
	}

	p.built = false
	return p
}

// SetCustomConditions updates the `CustomConditions` config values
func (p *Plan) SetCustomConditions(aliases map[string]CustomConditionFn, mode ...rune) *Plan {
	m := AppendMode
	if len(mode) > 0 && (mode[0] == OverwriteMode || mode[0] == WriteMode) {
		m = mode[0]
	}

	if m == OverwriteMode {
		p.config.CustomConditions = aliases
	} else {
		for key, val := range aliases {
			if _, ok := p.config.CustomConditions[key]; ok && m == AppendMode {
				continue
			}
			p.config.CustomConditions[key] = val
		}
	}

	p.built = false
	return p
}

// toCondition convert given interface to correct condition type
func toCondition(cond interface{}, vars []interface{}, not bool) (condition, error) {
	switch c := cond.(type) {
	case map[string]interface{}:
		return &mapConditions{value: c, not: not}, nil
	case []interface{}:
		if len(c) >= 2 {
			if cl, ok := c[0].(string); ok {
				return &rawConditions{clause: cl, vars: c[1:]}, nil
			}
		}
		return &orConditions{value: c, not: not}, nil
	case []map[string]interface{}:
		lenc := len(c)
		v := make([]interface{}, lenc)
		for i := 0; i < len(c); i++ {
			v[i] = c[i]
		}
		return &orConditions{value: v, not: not}, nil
	case string:
		return &rawConditions{clause: c, vars: vars, not: not}, nil
	default:
		return nil, &InvalidCond{cond: cond, vars: vars}
	}
}
