package gowhere

import (
	"sort"
	"strings"
)

// andConditions represents the list of conditions which will be tied by "AND"
type andConditions struct {
	value []interface{}
	naked bool
	not   bool
}

// orConditions represents the list of conditions which will be tied by "OR"
type orConditions struct {
	value []interface{}
	not   bool
}

// mapConditions represents the conditions in type of map
type mapConditions struct {
	value map[string]interface{}
	not   bool
}

// rawConditions represents a raw SQL condition
type rawConditions struct {
	clause string
	vars   []interface{}
	not    bool
}

// condition represents the condition interface
type condition interface {
	build(cfg *Config) (string, []interface{})
}

// CustomConditionFn represents the func signature which provide full access on the condition generating.
// Return value should be in form of condition, i.e a map or slice
// Return nil will exclude the condition from result
type CustomConditionFn func(key string, val interface{}, cfg *Config) interface{}

func (ac *andConditions) build(cfg *Config) (string, []interface{}) {
	sqls, vars := listBuild(ac.value, cfg)
	sql := strings.Join(sqls, " AND ")
	if sql != "" {
		if !ac.naked || ac.not {
			sql = "(" + sql + ")"
			if ac.not {
				sql = "NOT " + sql
			}
		}
	}

	return sql, vars
}

func (oc *orConditions) build(cfg *Config) (string, []interface{}) {
	sqls, vars := listBuild(oc.value, cfg)
	sql := strings.Join(sqls, " OR ")
	if sql != "" {
		sql = "(" + sql + ")"
		if oc.not {
			sql = "NOT " + sql
		}
	}

	return sql, vars
}

func (rc *rawConditions) build(cfg *Config) (string, []interface{}) {
	sql := rc.clause
	if sql != "" {
		sql = "(" + sql + ")"
		if rc.not {
			sql = "NOT " + sql
		}
	}
	return sql, rc.vars
}

func (mc *mapConditions) build(cfg *Config) (string, []interface{}) {
	vlen := len(mc.value)
	sqls := make([]string, 0, vlen)
	vars := make([]interface{}, 0)

	processFunc := func(key string, val interface{}) {
		var _sql string
		var _vars []interface{}

		if customCondFn, ok := cfg.CustomConditions[key]; ok {
			rawCond := customCondFn(key, val, cfg)
			if rawCond == nil {
				return
			}
			cond, err := toCondition(rawCond, []interface{}{}, false)
			if err != nil {
				if cfg.Strict {
					panic(&InvalidCond{cond: rawCond})
				}
				return
			}
			_sql, _vars = cond.build(cfg)

		} else {
			res := strings.Split(key, cfg.Separator)
			column := processColumn(res[0], cfg)
			var operator *Operator
			if len(res) > 1 {
				operator = findOperatorByName(res[1])
			} else {
				operator = findOperatorByValue(val)
			}

			if operator == nil {
				if cfg.Strict {
					panic(&InvalidCond{cond: key, vars: val})
				}
				return
			}

			_sql, _vars = operator.Build(column, val, cfg)

		}

		if _sql != "" {
			sqls = append(sqls, _sql)
			vars = append(vars, _vars...)
		}
	}

	if cfg.sort {
		keys := make([]string, 0, vlen)
		for k := range mc.value {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			val := mc.value[key]
			processFunc(key, val)
		}
	} else {
		for key, val := range mc.value {
			processFunc(key, val)
		}
	}

	sql := strings.Join(sqls, " AND ")
	if sql != "" {
		sql = "(" + sql + ")"
		if mc.not {
			sql = "NOT " + sql
		}
	}
	return sql, vars
}

// listBuild is shared func for building andConditions & orConditions
func listBuild(conds []interface{}, cfg *Config) ([]string, []interface{}) {
	sqls := make([]string, 0, len(conds))
	vars := make([]interface{}, 0)

	for i := 0; i < len(conds); i++ {
		var _sql string
		var _vars []interface{}

		switch c := conds[i].(type) {
		case map[string]interface{}:
			mconds := &mapConditions{value: c}
			_sql, _vars = mconds.build(cfg)
		case []interface{}:
			if len(c) > 0 {
				// test if it's in form of rawConditions
				if len(c) >= 2 {
					cl, cok := c[0].(string)
					va := c[1:]
					if cok {
						sconds := &rawConditions{clause: cl, vars: va}
						_sql, _vars = sconds.build(cfg)
					}
				}
				// it's not a rawConditions, consider as orConditions
				if _sql == "" {
					oconds := &orConditions{value: c}
					_sql, _vars = oconds.build(cfg)
				}
			}
		case *andConditions:
			_sql, _vars = c.build(cfg)
		case *orConditions:
			_sql, _vars = c.build(cfg)
		case *mapConditions:
			_sql, _vars = c.build(cfg)
		case *rawConditions:
			_sql, _vars = c.build(cfg)
		default:
			if cfg.Strict {
				panic(&InvalidCond{cond: c})
			}
			continue
		}

		if _sql != "" {
			sqls = append(sqls, _sql)
			vars = append(vars, _vars...)
		}
	}

	return sqls, vars
}

func processColumn(col string, cfg *Config) string {
	if alias, ok := cfg.ColumnAliases[col]; ok {
		col = alias
	}
	col = cfg.Dialect.QuoteIdentifier(col)

	if !strings.Contains(col, ".") && cfg.Table != "" {
		table := cfg.Dialect.QuoteIdentifier(cfg.Table)
		col = table + "." + col
	}

	return col
}
