package orm

import (
	"reflect"
	"testing"
)

// - Struct转Map，通过tags，可以使用自定义的tags

// go get -u github.com/wolfogre/gtag/cmd/gtag

type User struct {
	Id    int    `bson:"_id" orm:"id"`
	Name  string `bson:"name" orm:"name"`
	Email string `bson:"email"`
}

func TestTransfer(t *testing.T) {
	type args struct {
		tag    string
		object interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "tags parse",
			args: args{
				tag: "bson",
				object: User{
					Id:    1,
					Email: "xx",
					Name:  "ZZ",
				},
			},
			want: map[string]interface{}{
				"_id":   1,
				"email": "xx",
				"name":  "ZZ",
			},
		},
		{
			name: "orm tags parse",
			args: args{
				tag: "orm",
				object: User{
					Id:    1,
					Email: "xx",
					Name:  "ZZ",
				},
			},
			want: map[string]interface{}{
				"id":    1,
				"email": "xx",
				"name":  "ZZ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Transfer(tt.args.tag, tt.args.object); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transfer() = %v, want %v", got, tt.want)
			}
		})
	}
}
