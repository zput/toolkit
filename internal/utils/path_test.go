package utils

import (
	"testing"
)

func TestHowManySlash(t *testing.T) {
	type args struct {
		path          string
		targetSegment string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				path:          "toolkit/internal/utils/fake/path.go",
				targetSegment: "internal",
			},
			want: "../../../../",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HowManySlash(tt.args.path, tt.args.targetSegment); got != tt.want {
				t.Errorf("HowManySlash() = %v, want %v", got, tt.want)
			}
		})
	}
}
