package services

import (
	"net/http"
)

type HealthArgs struct {
}

type HealthReply struct {
	Status string `json:"status"`
}

type HealthService int

func (t *HealthService) Status(r *http.Request, args *HealthArgs, result *HealthReply) error {
	//	log.Printf("Health Check Called\n")
	*result = HealthReply{"up"}
	return nil
}
