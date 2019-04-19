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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := WithConfig(Config{sort: true}).Where(tt.args.cond, tt.args.vars...)
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
