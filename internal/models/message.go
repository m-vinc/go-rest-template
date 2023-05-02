package models

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	*kafka.Message
}

type EventType string

type GenericMessage struct {
	Event EventType `json:"event"`
}

type OperationMessage struct {
	Event EventType `json:"event"`
	At    time.Time `json:"at"`

	Metadata map[string]any `json:"metadata"`
}

var (
	OperationStartedEvent  EventType = "operations:started"
	OperationProgressEvent EventType = "operations:progress"
	OperationsEndedEvent   EventType = "operations:end"
)
