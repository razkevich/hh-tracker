package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Operations is an interface that can be used to mock out calls to the MongoDB driver
type Operations interface {
	GetSession() (Session, error)
	Collection() string
	WithSession(context.Context, Session, func(mongo.SessionContext) error) error
	ReplaceOne(context.Context, interface{}, interface{}) (*mongo.UpdateResult, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(context.Context, []interface{}, ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	FindOne(context.Context, interface{}, interface{}, ...*options.FindOneOptions) error
	FindOneAndUpdate(context.Context, interface{}, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) error
	Find(context.Context, interface{}, interface{}, ...*options.FindOptions) error
	FindBatch(context.Context, interface{}, ...*options.FindOptions) (Cursor, error)
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (Cursor, error)
	Delete(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
}

// Session is an interface around mongo.Session so we can mock out calls in unit tests
type Session interface {
	StartTransaction(...*options.TransactionOptions) error
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
	WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error),
		opts ...*options.TransactionOptions) (interface{}, error)
	EndSession(context.Context)
}

// Cursor is an interface wrapper around mongo.Cursor
type Cursor interface {
	All(ctx context.Context, results interface{}) error
	Next(ctx context.Context) bool
	Decode(val interface{}) error
	Close(ctx context.Context) error
	RemainingBatchLength() int
}

// NewOperations is a constructor for operations
func NewOperations(collectionName string, database *mongo.Database) Operations {
	return &operations{
		c: database.Collection(collectionName),
	}
}

type operations struct {
	c *mongo.Collection
}

func (oc *operations) Collection() string {
	return oc.c.Name()
}

func (oc *operations) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {
	return oc.c.CountDocuments(ctx, filter, opts...)
}

func (oc *operations) GetSession() (Session, error) {
	return oc.c.Database().Client().StartSession()
}

func (oc *operations) WithSession(ctx context.Context, sess Session, transaction func(sc mongo.SessionContext) error) error {
	return mongo.WithSession(ctx, sess.(mongo.Session), transaction)
}

func (oc *operations) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (Cursor, error) {
	return oc.c.Aggregate(ctx, pipeline, opts...)
}

func (oc *operations) InsertOne(ctx context.Context, doc interface{}, options ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return oc.c.InsertOne(ctx, doc, options...)
}

func (oc *operations) InsertMany(ctx context.Context, doc []interface{}, options ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return oc.c.InsertMany(ctx, doc, options...)
}

func (oc *operations) ReplaceOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return oc.c.ReplaceOne(ctx, filter, update)
}

func (oc *operations) UpdateOne(ctx context.Context, filter interface{}, update interface{}, options ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return oc.c.UpdateOne(ctx, filter, update, options...)
}

func (oc *operations) UpdateMany(ctx context.Context, filter interface{}, update interface{}, options ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return oc.c.UpdateMany(ctx, filter, update, options...)
}

func (oc *operations) FindOne(ctx context.Context, filter interface{}, result interface{}, options ...*options.FindOneOptions) error {
	return oc.c.FindOne(ctx, filter, options...).Decode(result)
}

func (oc *operations) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, result interface{}, options ...*options.FindOneAndUpdateOptions) error {
	return oc.c.FindOneAndUpdate(ctx, filter, update, options...).Decode(result)
}

func (oc *operations) Find(ctx context.Context, filter interface{}, result interface{}, options ...*options.FindOptions) error {
	cursor, err := oc.c.Find(ctx, filter, options...)
	if err != nil {
		return err
	}
	return cursor.All(ctx, result)
}

func (oc *operations) FindBatch(ctx context.Context, filter interface{}, options ...*options.FindOptions) (Cursor, error) {
	cursor, err := oc.c.Find(ctx, filter, options...)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func (oc *operations) Delete(ctx context.Context, filter interface{}, options ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return oc.c.DeleteOne(ctx, filter, options...)
}

func (oc *operations) DeleteMany(ctx context.Context, filter interface{}, options ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return oc.c.DeleteMany(ctx, filter, options...)
}
