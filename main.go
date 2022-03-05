package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collName = "test"
)

var (
	connString = os.Getenv("CONNSTRING")
	collection *mongo.Collection
)

type item struct {
	ID primitive.ObjectID `bson:"_id"`
	Name    string
	Number  int
	Content string
}

var (
	cache   = map[string]interface{}{}
	ctx     = context.Background()
	mongoDB *mongo.Database
)

func init() {
	conn, err := mongo.NewClient(options.Client().ApplyURI(connString))
	if err != nil {
		panic("mongo connection " + err.Error())
	}
	if err := conn.Connect(ctx); err != nil {
		panic(err)
	}

	mongoDB = conn.Database("test")
	collection = mongoDB.Collection(collName)

	// cache init
	items, err := GetItems()
	if err != nil {
		panic(err)
	}
	for _, i := range items {
		cache[i.ID.Hex()] = i
	}
}

func main() {

	stream, err := Stream()
	if err != nil {
		panic(err)
	}
	log.Println("db stream connected")

	for stream.Next(ctx) {
		var receiver struct {
			OperationType string
			FullDocument  item
			DocumentKey   struct {
				ID primitive.ObjectID `bson:"_id"`
			}
		}

		err = stream.Decode(&receiver)
		if err != nil {
			log.Println("decode error:", err)
			continue
		}
		log.Printf("%+v\n", receiver)

		// Update cache
		switch receiver.OperationType {
		case "delete":
			delete(cache, receiver.DocumentKey.ID.Hex())
		case "insert", "update":
			cache[receiver.DocumentKey.ID.Hex()] = receiver.FullDocument
		}

		log.Printf("cache: %+v\n", cache)
	}
}

func Stream() (*mongo.ChangeStream, error) {
	opts := options.ChangeStream()
	opts.SetFullDocument(options.UpdateLookup)

	pipe := mongo.Pipeline{}
	// pipe := mongo.Pipeline{{{"$match", bson.D{{"operationType", "insert"}}}}}

	return collection.Watch(ctx, pipe, opts)
}

func GetItems() ([]item, error) {
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var items []item
	if err = cur.All(ctx, &items); err != nil {
		return items, err
	}
	return items, nil
}

