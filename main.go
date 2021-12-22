package main

import (
	"context"
	"encoding/json"
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
	return &repository{coll: client.Database(db).Collection(col)}
}

type repository struct {
	coll *mongo.Collection
}

func (r *repository) createBook(doc interface{}) (*mongo.InsertOneResult, error) {
	result, err := r.coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) readBook(filter interface{}) ([]byte, error) {
	var result bson.M
	err := r.coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (r *repository) updateBook(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	result, err := r.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) deleteBook(filter interface{}) (*mongo.DeleteResult, error) {
	result, err := r.coll.DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Number of documents deleted: %d\n", result.DeletedCount)
	return result, err
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
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

	filter := bson.M{"_id": result.InsertedID}
	jsonData, err := repo.readBook(filter)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)

	update := bson.M{"$set": Book{
		Title:  "Bumi manusia",
		Author: "Pramoedya Ananta Toer",
		Year:   1980,
	}}
	updateResult, err := repo.updateBook(filter, update)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Documents matched: %v\n", updateResult.MatchedCount)
	fmt.Printf("Documents updated: %v\n", updateResult.ModifiedCount)

	jsonData, err = repo.readBook(filter)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)

	deleteResult, err := repo.deleteBook(filter)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Number of documents deleted: %d\n", deleteResult.DeletedCount)
}
