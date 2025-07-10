package internalhttp

import (
	"time"

	"github.com/hilltracer/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type createOrUpdateRequest struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	StartTime    time.Time     `json:"startTime"`
	Duration     time.Duration `json:"duration"`
	Description  string        `json:"description,omitempty"`
	UserID       string        `json:"userId"`
	NotifyBefore time.Duration `json:"notifyBefore,omitempty"`
}

type listResponse struct {
	Events []storage.Event `json:"events"`
}
