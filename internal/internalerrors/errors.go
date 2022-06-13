package internalerrors

import "errors"

var (
	// ErrExists is returned when an operation failed because an entity did not exist.
	ErrExists = errors.New("exists")
	// ErrConflict is returned when an operation failed because it would conflict with existing data.
	ErrConflict = errors.New("conflict")
	// ErrNotFound is returned when an operation failed because an entity already existed.
	ErrNotFound = errors.New("not found")
	// ErrStoreNotFound is returned when Store not found
	ErrStoreNotFound = errors.New("store not found")
	// ErrForbidden is returned when an operation failed because it was forbidden
	ErrForbidden = errors.New("forbidden")
	// ErrBadRequest is returned when an operation failed because of bad request
	ErrBadRequest = errors.New("bad request")
	// ErrUnprocessableEntity the syntax is correct however could not process it
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	// ErrMethodNotAllowed is returned when method not allowed
	ErrMethodNotAllowed = errors.New("method not allowed")
	// ErrStoreIDMismatch is returned when store id in header is not matched with id in path
	ErrStoreIDMismatch = errors.New("store id mismatch")
	// ErrNoValidRabbitMQHosts is returned when there's no valid RabbitMQ hosts
	ErrNoValidRabbitMQHosts = errors.New("no valid RabbitMQ hosts")
	// ErrOffsetLimitExceed is returned when provided offset exceeds the enforced limitation
	ErrOffsetLimitExceed = errors.New("offset limit has been exceeded")
	// ErrPageOffsetExceed is returned when provided offset exceeds the enforced limitation
	ErrPageOffsetExceed = errors.New("page offset has been exceeded enforced limitation")
	// ErrPageLimitExceed is returned when provided page limit exceeds the enforced limitation
	ErrPageLimitExceed = errors.New("page limit has been exceeded enforced limitation")
	// ErrSvcNotAvailable is returned when service is not available
	ErrSvcNotAvailable = errors.New("service not available")
	// ErrPublishTimeout is returned when publishing a message timeout
	ErrPublishTimeout = errors.New("timeout during publish")
	// ErrPublish is returned if there is an error while publishing
	ErrPublish = errors.New("failed to publish")
)
