package main

import (
	"context"
	"fmt"

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

func createBookRepository(uri, db, col string) *repository {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return &repository{collection: client.Database(db).Collection(col)}
}

type repository struct {
	collection *mongo.Collection
}

func (r *repository) createBook(book Book) (*Book, error) {
	res, err := r.collection.InsertOne(context.TODO(), book)
	if err != nil {
		return nil, err
	}

	book.ID = res.InsertedID.(primitive.ObjectID)
	return &book, nil
}

func (r *repository) readBook(id interface{}) (*Book, error) {
	filter := bson.M{"_id": id}
	var result bson.D
	err := r.collection.FindOne(context.TODO(), filter).Decode(&result)
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

func (r *repository) updateBook(id interface{}, book Book) (*Book, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": book}
	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *repository) deleteBook(id interface{}) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	return r.collection.DeleteMany(context.TODO(), filter)
}

func main() {
	uri := "mongodb+srv://admin:admin@cluster0.xtwwu.mongodb.net"
	database := "myDB"
	collection := "favorite_books"
	repo := createBookRepository(uri, database, collection)

	result, err := repo.createBook(Book{
		Title:  "Invisible Cities",
		Author: "Italo Calvino",
		Year:   1974,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted document with _id: %v\n", result.ID)

	book, err := repo.readBook(result.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", *book)

	updateResult, err := repo.updateBook(result.ID, Book{
		Title:  "Bumi manusia",
		Author: "Pramoedya Ananta Toer",
		Year:   1980,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Documents updated: %v\n", updateResult)

	book, err = repo.readBook(result.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", *book)

	deleteResult, err := repo.deleteBook(result.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Number of documents deleted: %d\n", deleteResult.DeletedCount)
}
