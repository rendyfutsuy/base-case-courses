package worker

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuybase-case-courses/constants"
)

// Job represents a generic task to be executed by a worker.
type Job struct {
	ID      uuid.UUID
	Type    constants.JobType
	Payload interface{} // This will hold structs from the 'payloads' package.
}

// JobQueue is the central channel for submitting background jobs.
var JobQueue chan Job
