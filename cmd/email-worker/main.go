package main

import (
	"log"

	"github.com/rendyfutsuybase-case-courses/modules/auth/tasks"
)

func main() {
	if err := tasks.RunEmailScheduler(); err != nil {
		log.Fatalf("email worker failed: %v", err)
	}
}
