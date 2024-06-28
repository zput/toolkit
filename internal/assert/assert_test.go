package assert

import (
	"github.com/shopspring/decimal"
	"testing"
)

func TestEqual(t *testing.T) {
	type args struct {
		cur interface{}
		exp interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test decimal",
			args: args{
				cur: GetDecimalPtr(decimal.RequireFromString("0.01")),
				exp: GetDecimalPtr(decimal.RequireFromString("0")),
			},
			want: false,
		},
		{
			name: "test decimal",
			args: args{
				cur: GetDecimalPtr(decimal.RequireFromString("0.01")),
				exp: GetDecimalPtr(decimal.RequireFromString("0.01")),
			},
			want: true,
		},
		{
			name: "test decimal",
			args: args{
				cur: GetDecimalPtr(decimal.RequireFromString("0.01")),
				exp: nil,
			},
			want: false,
		},
		{
			name: "test struct decimal",
			args: args{
				cur: DecimalStruct{
					First: GetDecimalPtr(decimal.RequireFromString("0.01")),
				},
				exp: DecimalStruct{
					First: GetDecimalPtr(decimal.RequireFromString("0.01")),
				},
			},
			want: true,
		},
		{
			name: "test struct decimal 2",
			args: args{
				cur: DecimalStruct{
					First: nil,
				},
				exp: DecimalStruct{
					First: GetDecimalPtr(decimal.RequireFromString("0.01")),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.args.cur, tt.args.exp); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func GetDecimalPtr(decimal decimal.Decimal) *decimal.Decimal {
	var tmp = decimal.Copy()
	return &tmp
}

type DecimalStruct struct {
	First *decimal.Decimal
}
