package service

import (
	"context"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/entity"
	log_entries "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/repository"
)

// ILogEntriesService is the interface describing service operations.
type ILogEntriesService interface {
	GetLogEntries(ctx context.Context, storeID string, offset, limit int, searchJSON string) (*entity.LogEntriesWithRelations, error)
}

// ServiceConfig represents the service configs
type ServiceConfig struct {
	// todo add some useful params here?
}

// LogEntriesService describes a service implementation.
type LogEntriesService struct {
	repository log_entries.Repository
	config     *ServiceConfig
}

// ProvideLogEntriesService is the wire function for LogEntry service.
func ProvideLogEntriesService(repository *log_entries.Repository, config *ServiceConfig) LogEntriesService {
	return LogEntriesService{
		repository: *repository,
		config:     config,
	}
}

// GetLogEntries lists all log entries. Responds with the list of log entries or a error.
func (a LogEntriesService) GetLogEntries(ctx context.Context, storeID string, offset, limit int, searchJSON string) (*entity.LogEntriesWithRelations, error) {

	return a.repository.Find(ctx)
}
