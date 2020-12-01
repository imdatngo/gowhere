package gowhere

import (
	"fmt"
	"reflect"
	"time"
)

// CustomBuildFn represents the function to build SQL string for the operator
type CustomBuildFn func(field string, value interface{}, cfg Config) (string, []interface{})

// ModValueFn represents the function to modify only the value before actually build the SQL
type ModValueFn func(value interface{}) interface{}

// Operator represents an alias for the SQL operator
type Operator struct {
	// Reference to an existing operator
	AliasOf string
	// The actual SQL operator. Default to "=" if empty
	Operator string
	// The SQL template. Default to "%s %s ?"
	// Note: This must contains 2 "%s" placeholders for the column & operator, and 1 "?" for the value
	Template string
	// The function to build the SQL condition in your own way. Ignored if AliasOf is provided, ignores Operator & Template.
	CustomBuild CustomBuildFn
	// Instead of customize the whole build func, you probably only want to modify the value a litle bit
	ModValue ModValueFn
}

// Build returns the SQL string & vars for a single condition.
func (o *Operator) Build(field string, value interface{}, cfg *Config) (string, []interface{}) {
	if o.CustomBuild != nil {
		return o.CustomBuild(field, value, *cfg)
	}

	operator := o.Operator
	if operator == "" {
		operator = "="
	}

	template := o.Template
	if template == "" {
		template = "%s %s ?"
	}

	if o.ModValue != nil {
		value = o.ModValue(value)
	}

	return fmt.Sprintf(template, field, operator), []interface{}{Utils.ToSQLVar(value)}
}

var (
	defaultOperator = &Operator{}

	// OperatorsList defines the list of built-in operators
	OperatorsList = map[string]*Operator{
		"exact":     defaultOperator,
		"iexact":    &Operator{Template: "LOWER(%s) %s LOWER(?)"},
		"notexact":  &Operator{Operator: "<>"},
		"inotexact": &Operator{Operator: "<>", Template: "LOWER(%s) %s LOWER(?)"},

		"gt":  &Operator{Operator: ">"},
		"lt":  &Operator{Operator: "<"},
		"gte": &Operator{Operator: ">="},
		"lte": &Operator{Operator: "<="},

		"startswith": &Operator{
			Operator: "LIKE",
			ModValue: func(value interface{}) interface{} {
				return Utils.ToString(value) + "%"
			},
		},
		"istartswith": &Operator{
			Operator: "LIKE",
			Template: "LOWER(%s) %s LOWER(?)",
			ModValue: func(value interface{}) interface{} {
				return Utils.ToString(value) + "%"
			},
		},
		"endswith": &Operator{
			Operator: "LIKE",
			ModValue: func(value interface{}) interface{} {
				return "%" + Utils.ToString(value)
			},
		},
		"iendswith": &Operator{
			Operator: "LIKE",
			Template: "LOWER(%s) %s LOWER(?)",
			ModValue: func(value interface{}) interface{} {
				return "%" + Utils.ToString(value)
			},
		},
		"contains": &Operator{
			Operator: "LIKE",
			ModValue: func(value interface{}) interface{} {
				return "%" + Utils.ToString(value) + "%"
			},
		},
		"icontains": &Operator{
			Operator: "LIKE",
			Template: "LOWER(%s) %s LOWER(?)",
			ModValue: func(value interface{}) interface{} {
				return "%" + Utils.ToString(value) + "%"
			},
		},
		"in": &Operator{
			Operator: "IN",
			Template: "%s %s (?)",
			ModValue: func(value interface{}) interface{} {
				return Utils.ToSlice(value)
			},
		},
		"date": &Operator{
			Template: "DATE(%s) %s ?",
			ModValue: func(value interface{}) interface{} {
				return Utils.ToDate(value)
			},
		},
		"between": &Operator{
			CustomBuild: func(field string, value interface{}, cfg Config) (string, []interface{}) {
				var from interface{}
				var to interface{}

				if vi, ok := value.([]interface{}); ok && len(vi) >= 2 {
					from = vi[0]
					to = vi[1]
				} else if vs, ok := value.([]string); ok && len(vs) >= 2 {
					from = vs[0]
					to = vs[1]
				} else if vt, ok := value.([]time.Time); ok && len(vt) >= 2 {
					from = vt[0]
					to = vt[1]
				} else {
					return "", []interface{}{}
				}

				return fmt.Sprintf("%s BETWEEN ? AND ?", field), []interface{}{Utils.ToDateTime(from), Utils.ToDateTime(to)}
			},
		},
		"isnull": &Operator{
			CustomBuild: func(field string, value interface{}, cfg Config) (string, []interface{}) {
				operator := "IS NULL"
				if null, ok := value.(bool); ok && !null {
					operator = "IS NOT NULL"
				}
				return fmt.Sprintf("%s %s", field, operator), []interface{}{}
			},
		},
		"datebetween": &Operator{
			CustomBuild: func(field string, value interface{}, cfg Config) (string, []interface{}) {
				var from interface{}
				var to interface{}

				if vi, ok := value.([]interface{}); ok && len(vi) >= 2 {
					from = vi[0]
					to = vi[1]
				} else if vs, ok := value.([]string); ok && len(vs) >= 2 {
					from = vs[0]
					to = vs[1]
				} else if vt, ok := value.([]time.Time); ok && len(vt) >= 2 {
					from = vt[0]
					to = vt[1]
				} else {
					return "", []interface{}{}
				}

				return fmt.Sprintf("DATE(%s) BETWEEN ? AND ?", field), []interface{}{Utils.ToDate(from), Utils.ToDate(to)}
			},
		},
	}
)

func findOperatorByName(name string) *Operator {
	if op, ok := OperatorsList[name]; ok {
		if op.AliasOf != "" {
			if alias := findOperatorByName(op.AliasOf); alias != nil {
				return alias
			}
		}
		return op
	}
	return nil
}

func findOperatorByValue(value interface{}) *Operator {
	if value == nil {
		return findOperatorByName("isnull")
	}

	vt := reflect.TypeOf(value)
	if vt.Kind() == reflect.Array || vt.Kind() == reflect.Slice {
		return findOperatorByName("in")
	}

	return defaultOperator
}
