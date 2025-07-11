package api

import (
	"time"
)

// Meta provides common fields for all API objects
type Meta struct {
	ID        string     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// List provides pagination for all list responses
type List struct {
	Kind  string      `json:"kind"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

// EventType defines the types of events in the system
type EventType string

const (
	CreateEventType EventType = "create"
	UpdateEventType EventType = "update"
	DeleteEventType EventType = "delete"
)

// Event represents a system event
type Event struct {
	Meta
	Source         string     `json:"source"`
	SourceID       string     `json:"source_id"`
	EventType      EventType  `json:"event_type"`
	ReconciledDate *time.Time `json:"reconciled_date,omitempty"`
}

// ListQuery represents query parameters for list operations
type ListQuery struct {
	Page    int
	Size    int
	Search  string
	OrderBy string
	Fields  string
}

// PaginationMeta provides pagination metadata
type PaginationMeta struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// ResourceReference provides a reference to any resource
type ResourceReference struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	Href string `json:"href"`
}
