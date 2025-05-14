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

/*
func (r *ExampleRepository) FindAll() (*[]Example, error) {
	collection := r.client.Client.Database(r.dbName).Collection(r.collection)

	var examples []Example
	var cursor *mongo.Cursor
	// cursor, err = collection.Find(context.Background(), {})
	(context.Background()).Decode(&examples)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("example not found")
		}
		return nil, err
	}

	return &example, nil
}*/

// GetAllExamples retrieves all examples from the collection
func (r *ExampleRepository) GetAllExamples() ([]Example, error) {
	collection := r.client.Client.Database(r.dbName).Collection(r.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find options (optional)
	findOptions := options.Find()
	// findOptions.SetLimit(100) // You can set limit if needed

	// Execute query
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var examples []Example
	if err = cursor.All(ctx, &examples); err != nil {
		return nil, err
	}

	return examples, nil
}

// GetExamplesWithFilter retrieves examples matching the given filter
func (r *ExampleRepository) GetExamplesWithFilter(filter bson.M) ([]Example, error) {
	collection := r.client.Client.Database(r.dbName).Collection(r.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var examples []Example
	if err = cursor.All(ctx, &examples); err != nil {
		return nil, err
	}

	return examples, nil
}

// GetPaginatedExamples retrieves paginated results
func (r *ExampleRepository) GetPaginatedExamples(page, limit int64) ([]Example, error) {
	collection := r.client.Client.Database(r.dbName).Collection(r.collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSkip((page - 1) * limit)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by created_at descending

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var examples []Example
	if err = cursor.All(ctx, &examples); err != nil {
		return nil, err
	}

	return examples, nil
}

// Add other repository methods (FindAll, Update, Delete, etc.)
