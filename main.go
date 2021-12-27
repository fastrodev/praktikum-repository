package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Year   int                `bson:"year_published,omitempty"`
}

func createBookRepository(
	collection *mongo.Collection,
	timeout time.Duration,
) *repository {
	return &repository{
		collection: collection,
		timeout:    timeout,
	}
}

func createCollection(
	ctx context.Context,
	timeout time.Duration,
	uri, db, col string,
) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client.Database(db).Collection(col), nil
}

type repository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func (r *repository) createBook(ctx context.Context, book Book) (*Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	res, err := r.collection.InsertOne(ctx, book)
	if err != nil {
		return nil, err
	}

	book.ID = res.InsertedID.(primitive.ObjectID)
	return &book, nil
}

func (r *repository) readBook(ctx context.Context, id interface{}) (*Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"_id": id}
	var result bson.D
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	docByte, err := bson.Marshal(result)
	if err != nil {
		return nil, err
	}

	var book Book
	err = bson.Unmarshal(docByte, &book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *repository) updateBook(ctx context.Context, id interface{}, book Book) (*Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": book}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *repository) deleteBook(ctx context.Context, id interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"_id": id}
	return r.collection.DeleteMany(ctx, filter)
}

func main() {
	uri := "mongodb+srv://admin:admin@cluster0.xtwwu.mongodb.net"
	database := "myDB"
	col := "favorite_books"
	ctx := context.Background()
	timeout := 10 * time.Second
	collection, err := createCollection(ctx, timeout, uri, database, col)
	if err != nil {
		panic(err)
	}
	repo := createBookRepository(collection, timeout)

	result, err := repo.createBook(
		ctx,
		Book{
			Title:  "Invisible Cities",
			Author: "Italo Calvino",
			Year:   1974,
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted document with _id: %v\n", result.ID)

	book, err := repo.readBook(ctx, result.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", *book)

	updateResult, err := repo.updateBook(
		ctx,
		result.ID,
		Book{
			Title:  "Bumi manusia",
			Author: "Pramoedya Ananta Toer",
			Year:   1980,
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Documents updated: %v\n", updateResult)

	book, err = repo.readBook(ctx, result.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", *book)

	deleteResult, err := repo.deleteBook(ctx, result.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Number of documents deleted: %d\n", deleteResult.DeletedCount)
}
