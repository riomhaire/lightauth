package services

import (
	"net/http"

	"github.com/thoas/stats"
)

type StatsArgs struct {
}

type StatsService int
type StatsReply struct {
	Application string      `json:"application"`
	Version     string      `json:"version"`
	Host        string      `json:"host"`
	Statistics  *stats.Data `json:"statistics"`
}

func (t *StatsService) Status(r *http.Request, args *StatsArgs, result *StatsReply) error {
	//	log.Printf("Statistics Check Called\n")
	reply := StatsReply{}
	reply.Application = "Authorization"
	reply.Host = GetLocalIP()
	reply.Version = Version
	reply.Statistics = StatsMiddleware.Data()

	*result = reply
	return nil
}
