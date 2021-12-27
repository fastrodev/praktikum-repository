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
	ctx context.Context,
	timeout time.Duration,
	uri, db, col string,
) *repository {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return &repository{
		collection: client.Database(db).Collection(col),
		timeout:    timeout,
	}
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
	collection := "favorite_books"
	ctx := context.Background()
	timeout := 10 * time.Second
	repo := createBookRepository(ctx, timeout, uri, database, collection)

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
