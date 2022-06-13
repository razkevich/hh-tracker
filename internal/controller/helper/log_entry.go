package helper

import (
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/entity"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/shared"
	"time"
)

// FullLogEntry includes relationships and meta
type FullLogEntry struct {
	ID      string             `json:"id,omitempty"`
	StoreID string             `json:"store_id,omitempty"`
	Name    string             `json:"name,omitempty"`
	Type    string             `json:"type,omitempty"`
	Rels    *LogEntryRelations `json:"relationships,omitempty"`
	Meta    Meta               `json:"meta,omitempty"`
	Links   *Links             `json:"links,omitempty"`
} // @name Response.FullLogEntry

// BuildLogEntryWithLinks converts log entry into the model with links
func BuildLogEntryWithLinks(logEntry entity.LogEntry, links Links) FullLogEntry {
	timestamps := Timestamps{
		CreatedAt: FormattedTimestamp{time.Now()},
	}
	fullLogEntry := FullLogEntry{
		ID:      logEntry.ID.String(),
		StoreID: logEntry.StoreID.String(),
		Name:    logEntry.Name,
		Type:    shared.LogEntry,
		Rels:    &LogEntryRelations{ResourcePath: LogEntryResourcePathRelation{URL: "dummy"}},
		Meta:    Meta{Timestamps: &timestamps},
		Links:   &links,
	}

	return fullLogEntry
}

// LogEntryRelations LogEntryRelations
type LogEntryRelations struct {
	ResourcePath LogEntryResourcePathRelation `json:"resource_path,omitempty"`
}

// LogEntryResourcePathRelation LogEntryResourcePathRelation
type LogEntryResourcePathRelation struct {
	URL string `json:"url,omitempty"`
}
