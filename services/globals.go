package services

import (
	"github.com/muesli/cache2go"
	"github.com/thoas/stats"
)

// Application Parameters
var SessionSecret *string
var SessionPeriod *int
var SessionsFile *string // Known sessions - like long term API keys
var UserFile *string     // Known users

// Contains the global variables used
var StatsMiddleware = stats.New()

// User stuff
var users = make(map[string]User)

// Session stuff
var knownSessions []string
var sessions = cache2go.Cache("sessionCache")
var ApplicationSessionService = new(SessionService)
