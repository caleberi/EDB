package repository

import (
	"context"
	"yc-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositories struct {
	UserRepository         Repository[models.User]
	EmployeeRepository     Repository[models.Employee]
	DisbursementRepository Repository[models.Disbursement]
}

func InitRepositories(db *mongo.Database) *Repositories {
	// register all collection here so we can provide via gin.context
	userRepo := NewRepository[models.User](db.Collection("users"))
	employeeRepo := NewRepository[models.Employee](db.Collection("employees"))
	disbursementRepo := NewRepository[models.Disbursement](db.Collection("disbursement"))
	return &Repositories{
		UserRepository:         userRepo,
		EmployeeRepository:     employeeRepo,
		DisbursementRepository: disbursementRepo,
	}
}

// IRepository defines the methods that a repository must implement.
type IRepository[T comparable] interface {
	Create(ctx context.Context, document T) (any, error)
	FindOneById(ctx context.Context, id primitive.ObjectID) (*T, error)
	FindOne(ctx context.Context, filter bson.D) (*T, error)
	FindMany(ctx context.Context, filter bson.D) ([]T, error)
	UpdateOneById(ctx context.Context, id primitive.ObjectID, document T) error
	UpdateMany(ctx context.Context, filter bson.D, document T) error
	DeleteById(ctx context.Context, id primitive.ObjectID) error
	DeleteMany(ctx context.Context, filter bson.D) error
	Count(ctx context.Context, filter bson.D) (int64, error)
	CreateIndex(ctx context.Context, keys bson.D, opt *options.IndexOptions) (string, error)
	EstimatedDocumentCount(ctx context.Context) (int64, error)
	Aggregate(ctx context.Context, pipeline mongo.Pipeline, opts ...*options.AggregateOptions) ([]*T, error)
}

// Repository is a MongoDB repository implementation.
type Repository[T comparable] struct {
	collection *mongo.Collection // MongoDB collection
}

// NewRepository creates a new instance of Repository.
func NewRepository[T comparable](collection *mongo.Collection) Repository[T] {
	return Repository[T]{collection: collection}
}

// Create inserts a document into the MongoDB collection.
func (r *Repository[T]) Create(ctx context.Context, document T) (any, error) {
	result, err := r.collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

// FindOneById finds a single document by its ID in the MongoDB collection.
func (r *Repository[T]) FindOneById(ctx context.Context, id primitive.ObjectID) (*T, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	var result T
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindOne finds a single document based on the provided filter in the MongoDB collection.
func (r *Repository[T]) FindOne(ctx context.Context, filter bson.D) (*T, error) {
	var result T
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindMany finds multiple documents based on the provided filter in the MongoDB collection.
func (r *Repository[T]) FindMany(ctx context.Context, filter bson.D) ([]T, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var result T
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// UpdateOneById updates a single document by its ID in the MongoDB collection.
func (r *Repository[T]) UpdateOneById(ctx context.Context, id primitive.ObjectID, document T) error {
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: document}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateMany updates multiple documents based on the provided filter in the MongoDB collection.
func (r *Repository[T]) UpdateMany(ctx context.Context, filter bson.D, document T) error {
	update := bson.D{{Key: "$set", Value: document}}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

// DeleteById deletes a single document by its ID from the MongoDB collection.
func (r *Repository[T]) DeleteById(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// DeleteMany deletes multiple documents based on the provided filter from the MongoDB collection.
func (r *Repository[T]) DeleteMany(ctx context.Context, filter bson.D) error {
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// Count returns the number of documents that match the given filter in the MongoDB collection.
func (r *Repository[T]) Count(ctx context.Context, filter bson.D) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateIndex creates an index in the MongoDB collection based on the specified keys and options.
func (r *Repository[T]) CreateIndex(ctx context.Context, keys bson.D, opt *options.IndexOptions) (string, error) {
	index := mongo.IndexModel{
		Keys:    keys,
		Options: opt,
	}
	return r.collection.Indexes().CreateOne(ctx, index)
}

// EstimatedDocumentCount returns an estimate of the number of documents in the MongoDB collection.
func (r *Repository[T]) EstimatedDocumentCount(ctx context.Context) (int64, error) {
	count, err := r.collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Aggregate performs an aggregation operation on the MongoDB collection based on the provided pipeline and options.
func (r *Repository[T]) Aggregate(ctx context.Context, pipeline mongo.Pipeline, opts ...*options.AggregateOptions) ([]*T, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*T
	for cursor.Next(ctx) {
		var result T
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
