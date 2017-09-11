package services

import (
	"net/http"

	"github.com/thoas/stats"
)

type StatsArgs struct {
}

type StatsService int

func (t *StatsService) Status(r *http.Request, args *StatsArgs, result *stats.Data) error {
	//	log.Printf("Statistics Check Called\n")
	*result = *StatsMiddleware.Data()
	return nil
}
