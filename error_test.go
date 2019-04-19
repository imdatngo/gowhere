package gowhere

import "testing"

func TestInvalidCond_Error(t *testing.T) {
	type fields struct {
		cond interface{}
		vars interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "invalid condition",
			fields: fields{
				cond: []string{"Hello ?"},
				vars: []interface{}{"World!"},
			},
			want: "Invalid Conditions: [Hello ?] [World!]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &InvalidCond{
				cond: tt.fields.cond,
				vars: tt.fields.vars,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("InvalidCond.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
