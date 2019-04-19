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
	condition := p.toCondition(cond, vars, false)
	if condition == nil {
		return p
	}

	p.conditions.value = append(p.conditions.value, condition)
	p.built = false

	return p
}

// Or wraps all current conditions and ties with the new "cond" by OR operator
func (p *Plan) Or(cond interface{}, vars ...interface{}) *Plan {
	condition := p.toCondition(cond, vars, false)
	if condition == nil {
		return p
	}

	p.conditions.naked = false
	p.conditions = &andConditions{naked: true, value: []interface{}{&orConditions{value: []interface{}{p.conditions, condition}}}}
	p.built = false

	return p
}

// Not works similar to Where but reverses the condition operator(s)
func (p *Plan) Not(cond interface{}, vars ...interface{}) *Plan {
	condition := p.toCondition(cond, vars, true)
	if condition == nil {
		return p
	}

	p.conditions.value = append(p.conditions.value, condition)
	p.built = false

	return p
}

// Build builds the SQL clause and vars with given conditions
func (p *Plan) Build() *Plan {
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

// // SetTable set the `Table` config value
// func (p *Plan) SetTable(value string) *Plan {
// 	p.config.Table = value
// 	return p
// }

// // SetColumnAliases set the `Table` config value
// func (p *Plan) SetColumnAliases(aliases map[string]string, mode ...rune) *Plan {
// 	m := 'a'
// 	if len(mode) > 0 && (mode[0] == 'o' || mode[0] == 'w') {
// 		m = mode[0]
// 	}

// 	if m == 'o' {
// 		p.config.ColumnAliases = aliases
// 	} else {
// 		for key, val := range aliases {
// 			if _, ok := p.config.ColumnAliases[key]; ok && m == 'a' {
// 				continue
// 			}
// 			p.config.ColumnAliases[key] = val
// 		}
// 	}

// 	return p
// }

// // se

// toCondition convert given interface to correct condition type
func (p *Plan) toCondition(cond interface{}, vars []interface{}, not bool) interface{} {
	switch c := cond.(type) {
	case map[string]interface{}:
		return &mapConditions{value: c, not: not}
	case []interface{}:
		return &orConditions{value: c, not: not}
	case []map[string]interface{}:
		lenc := len(c)
		v := make([]interface{}, lenc)
		for i := 0; i < len(c); i++ {
			v[i] = c[i]
		}
		return &orConditions{value: v, not: not}
	case string:
		return &rawConditions{clause: c, vars: vars, not: not}
	default:
		if p.config.Strict {
			p.Error = &InvalidCond{cond: cond, vars: vars}
		}
		return nil
	}
}
