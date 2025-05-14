package repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient struct {
	Client *mongo.Client
}

func NewMongoDBClient(uri string) (*MongoDBClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &MongoDBClient{Client: client}, nil
}

func (c *MongoDBClient) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c.Client.Disconnect(ctx)
}

type ExampleRepository struct {
	client     *MongoDBClient
	dbName     string
	collection string
}

func NewExampleRepository(client *MongoDBClient, dbName string) *ExampleRepository {
	return &ExampleRepository{
		client:     client,
		dbName:     dbName,
		collection: "examples",
	}
}

// Example model for MongoDB
type Example struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (r *ExampleRepository) Create(example *Example) (*Example, error) {
	collection := r.client.Client.Database(r.dbName).Collection(r.collection)

	example.CreatedAt = time.Now()
	example.UpdatedAt = time.Now()

	res, err := collection.InsertOne(context.Background(), example)
	if err != nil {
		return nil, err
	}

	example.ID = res.InsertedID.(primitive.ObjectID)
	return example, nil
}

func (r *ExampleRepository) FindByID(id string) (*Example, error) {
	collection := r.client.Client.Database(r.dbName).Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var example Example
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&example)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("example not found")
		}
		return nil, err
	}

	return &example, nil
}

// Add other repository methods (FindAll, Update, Delete, etc.)
