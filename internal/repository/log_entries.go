package logentries

import (
	"context"
	"github.com/google/uuid"
	mongoDriver "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/driver"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/entity"
	"time"

	mg "go.mongodb.org/mongo-driver/mongo"
)

// Repository represents a data repository.
type Repository struct {
	operations mongoDriver.Operations
	timeout    time.Duration
}

// NewRepository is a constructor function to create a new Repository
func NewRepository(operations mongoDriver.Operations, timeout time.Duration) *Repository {
	return &Repository{
		operations: operations,
		timeout:    timeout,
	}
}

// ProvideRepository factory method for wire
func ProvideRepository(config *mongoDriver.MongoConfig, database *mg.Database) *Repository {
	operations := mongoDriver.NewOperations("log_entries", database)
	return NewRepository(operations, time.Duration(config.Timeout)*time.Millisecond)
}

// Find returns a single product with a specific ID
func (r Repository) Find(ctx context.Context) (*entity.LogEntriesWithRelations, error) {
	return &entity.LogEntriesWithRelations{LogEntries: []*entity.LogEntry{{
		ID:      uuid.New(),
		StoreID: uuid.New(),
		Name:    "world",
	},
	}}, nil
}
