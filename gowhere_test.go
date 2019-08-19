package gowhere

import (
	"reflect"
	"testing"
)

func TestWhere(t *testing.T) {
	type args struct {
		cond interface{}
		vars []interface{}
	}
	tests := []struct {
		name     string
		cfg      Config
		args     args
		wantSQL  string
		wantVars []interface{}
	}{
		{
			name: "map input",
			args: args{
				cond: map[string]interface{}{
					"budget__gte": 1000,
					"name":        "Gopher",
				},
			},
			wantSQL:  `("budget" >= ? AND "name" = ?)`,
			wantVars: []interface{}{1000, "Gopher"},
		},
		{
			name: "string input",
			args: args{
				cond: "name = ? and budget >= ?",
				vars: []interface{}{"Go", 2000},
			},
			wantSQL:  "(name = ? and budget >= ?)",
			wantVars: []interface{}{"Go", 2000},
		},
		{
			name: "slice input",
			args: args{
				cond: []interface{}{
					map[string]interface{}{
						"name": "Gopher",
					},
					map[string]interface{}{
						"budget__gte": 3000,
					},
					[]interface{}{"date between ? and ?", "2019-04-17", "2019-04-18"},
				},
			},
			wantSQL:  `(("name" = ?) OR ("budget" >= ?) OR (date between ? and ?))`,
			wantVars: []interface{}{"Gopher", 3000, "2019-04-17", "2019-04-18"},
		},
		{
			name: "column alias",
			cfg: Config{
				ColumnAliases: map[string]string{
					"name":   "full_name",
					"budget": "price",
				},
			},
			args: args{
				cond: map[string]interface{}{
					"budget__gte": 2000,
					"name":        "Go",
				},
			},
			wantSQL:  `("price" >= ? AND "full_name" = ?)`,
			wantVars: []interface{}{2000, "Go"},
		},
		{
			name: "custom condition with map",
			cfg: Config{
				CustomConditions: map[string]CustomConditionFn{
					"search": func(key string, val interface{}, cfg *Config) interface{} {
						return []interface{}{
							map[string]interface{}{"first_name__contains": val},
							map[string]interface{}{"last_name__contains": val},
						}
					},
				},
			},
			args: args{
				cond: map[string]interface{}{
					"budget__gte": 2000,
					"search":      "Go",
				},
			},
			wantSQL:  `("budget" >= ? AND (("first_name" LIKE ?) OR ("last_name" LIKE ?)))`,
			wantVars: []interface{}{2000, "%Go%", "%Go%"},
		},
		{
			name: "custom condition with raw sql",
			cfg: Config{
				CustomConditions: map[string]CustomConditionFn{
					"search": func(key string, val interface{}, cfg *Config) interface{} {
						val = "%" + val.(string) + "%"
						return []interface{}{"first_name like ? or last_name like ?", val, val}
					},
				},
			},
			args: args{
				cond: map[string]interface{}{
					"budget__gte": 2000,
					"search":      "Go",
				},
			},
			wantSQL:  `("budget" >= ? AND (first_name like ? or last_name like ?))`,
			wantVars: []interface{}{2000, "%Go%", "%Go%"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cfg.sort = true
			plan := WithConfig(tt.cfg).Where(tt.args.cond, tt.args.vars...)
			sql := plan.SQL()
			vars := plan.Vars()

			if plan.Error != nil {
				t.Errorf("unexpected error: %+v", plan.Error)
			}

			if sql != tt.wantSQL {
				t.Errorf("sql = %v, want %v", sql, tt.wantSQL)
			}

			if !reflect.DeepEqual(vars, tt.wantVars) {
				t.Errorf("vars = %v, want %v", vars, tt.wantVars)
			}
		})
	}
}
