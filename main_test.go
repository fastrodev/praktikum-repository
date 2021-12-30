package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	uri        string = "mongodb+srv://admin:admin@cluster0.xtwwu.mongodb.net"
	database   string = "myDB"
	col        string = "favorite_books"
	id         primitive.ObjectID
	ctx        context.Context = context.Background()
	timeout    time.Duration   = 10 * time.Second
	collection *mongo.Collection
)

func Test_createBookRepository(t *testing.T) {
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

func Test_repository_createBook(t *testing.T) {
	collection, _ = createCollection(ctx, timeout, uri, database, col)
	type fields struct {
		collection *mongo.Collection
		timeout    time.Duration
	}
	type args struct {
		ctx  context.Context
		book Book
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Book
		wantErr bool
	}{

		{
			name: "empty book",
			fields: fields{
				collection: collection,
				timeout:    timeout,
			},
			args: args{
				ctx: ctx,
			},
			want:    &Book{},
			wantErr: false,
		},
		{
			name: "success create",
			fields: fields{
				collection: collection,
				timeout:    timeout,
			},
			args: args{
				ctx: ctx,
				book: Book{
					Title:  "Invisible Cities",
					Author: "Italo Calvino",
					Year:   1974,
				},
			},
			want: &Book{
				Title:  "Invisible Cities",
				Author: "Italo Calvino",
				Year:   1974,
			},
			wantErr: false,
		},
		{
			name: "context deadline exceeded",
			fields: fields{
				collection: collection,
			},
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				collection: tt.fields.collection,
				timeout:    tt.fields.timeout,
			}
			got, err := r.createBook(tt.args.ctx, tt.args.book)
			actualErr := err != nil
			require.Equal(t, tt.wantErr, actualErr)
			if got != nil {
				id = got.ID
				require.Equal(t, tt.want.Author, got.Author)
				require.Equal(t, tt.want.Title, got.Title)
				require.Equal(t, tt.want.Year, got.Year)
			}
		})
	}
}

func Test_repository_readBook(t *testing.T) {
	type fields struct {
		collection *mongo.Collection
		timeout    time.Duration
	}
	type args struct {
		ctx context.Context
		id  interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Book
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				collection: collection,
				timeout:    timeout,
			},
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: &Book{
				Title:  "Invisible Cities",
				Author: "Italo Calvino",
				Year:   1974,
			},
			wantErr: false,
		},
		{
			name: "fail",
			fields: fields{
				collection: collection,
				timeout:    timeout,
			},
			args: args{
				ctx: ctx,
				id:  nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				collection: tt.fields.collection,
				timeout:    tt.fields.timeout,
			}
			got, err := r.readBook(tt.args.ctx, tt.args.id)
			actualErr := (err != nil)
			require.Equal(t, tt.wantErr, actualErr)
			if got != nil {
				require.Equal(t, tt.want.Author, got.Author)
				require.Equal(t, tt.want.Title, got.Title)
				require.Equal(t, tt.want.Year, got.Year)
			}
		})
	}
}

func Test_repository_updateBook(t *testing.T) {
	type fields struct {
		collection *mongo.Collection
		timeout    time.Duration
	}
	type args struct {
		ctx  context.Context
		id   interface{}
		book Book
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Book
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				collection: collection,
				timeout:    timeout,
			},
			args: args{
				ctx: ctx,
				id:  id,
				book: Book{
					Title:  "Bumi manusia",
					Author: "Pramoedya Ananta Toer",
					Year:   1980,
				},
			},
			want: &Book{
				Title:  "Bumi manusia",
				Author: "Pramoedya Ananta Toer",
				Year:   1980,
			},
			wantErr: false,
		},
		{
			name: "fail",
			fields: fields{
				collection: collection,
			},
			args: args{
				ctx: ctx,
				id:  id,
				book: Book{
					Title:  "Bumi manusia",
					Author: "Pramoedya Ananta Toer",
					Year:   1980,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				collection: tt.fields.collection,
				timeout:    tt.fields.timeout,
			}
			got, err := r.updateBook(tt.args.ctx, tt.args.id, tt.args.book)
			actualErr := (err != nil)
			require.Equal(t, tt.wantErr, actualErr)
			if got != nil {
				require.Equal(t, tt.want.Author, got.Author)
				require.Equal(t, tt.want.Title, got.Title)
				require.Equal(t, tt.want.Year, got.Year)
			}
		})
	}
}

func Test_repository_deleteBook(t *testing.T) {
	type fields struct {
		collection *mongo.Collection
		timeout    time.Duration
	}
	type args struct {
		ctx context.Context
		id  interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mongo.DeleteResult
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				collection: collection,
				timeout:    timeout,
			},
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: &mongo.DeleteResult{
				DeletedCount: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				collection: tt.fields.collection,
				timeout:    tt.fields.timeout,
			}
			got, err := r.deleteBook(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.deleteBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repository.deleteBook() = %v, want %v", got, tt.want)
			}
		})
	}
}
