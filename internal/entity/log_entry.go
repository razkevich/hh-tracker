package entity

import (
	"fmt"
	"github.com/google/uuid"
)

// LogEntry the base type in the AM domain.
type LogEntry struct {
	ID      uuid.UUID `gorm:"column:id" json:"id,omitempty" form:"id"`
	StoreID uuid.UUID `gorm:"column:store_id" json:"-" format:"uuid"`
	Name    string    `gorm:"column:name" json:"name"`
}

// LogEntries the list of log entries and total number of log entries
type LogEntries struct {
	LogEntries []LogEntry
	TotalCount int64
}

// GetSelfLink function to get self link for LogEntry
func (a LogEntry) GetSelfLink(baseURL string) string {
	return fmt.Sprintf("%s/v2/personal-data/logs/%s", baseURL, a.ID)
}

// LogEntriesWithRelations the list of Log Entries with parent and ancestors and total number of Log Entries
type LogEntriesWithRelations struct {
	LogEntries []*LogEntry
	TotalCount int64
}
