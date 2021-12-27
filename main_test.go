package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func Test_createBookRepository(t *testing.T) {
	uri := "mongodb+srv://admin:admin@cluster0.xtwwu.mongodb.net"
	database := "myDB"
	col := "favorite_books"
	ctx := context.Background()
	timeout := 10 * time.Second
	successCollection, _ := createCollection(ctx, timeout, uri, database, col)
	errCollection, _ := createCollection(ctx, timeout, "", database, col)

	type args struct {
		collection *mongo.Collection
		timeout    time.Duration
	}
	tests := []struct {
		name string
		args args
		want *repository
	}{
		{
			name: "success create repository",
			args: args{
				collection: successCollection,
				timeout:    timeout,
			},
			want: &repository{
				collection: successCollection,
				timeout:    timeout,
			},
		},
		{
			name: "fail create collection",
			args: args{
				collection: errCollection,
				timeout:    timeout,
			},
			want: &repository{
				collection: errCollection,
				timeout:    timeout,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createBookRepository(tt.args.collection, tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createBookRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
