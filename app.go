package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/riomhaire/lightauth/services"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
	"github.com/urfave/negroni"
)

var (
	version string = "0.4"
)

func main() {
	// Command lines
	log.Printf("%s version %s\n", os.Args[0], version)

	services.SessionSecret = flag.String("sessionSecret", "secret", "Master key which is used to generate system jwt")
	services.SessionPeriod = flag.Int("sessionPeriod", 3600, "How many seconds before sessions expires")
	services.UserFile = flag.String("usersFile", "users.csv", "List of Users and salted/hashed password with their roles")
	services.SessionsFile = flag.String("sessionFile", "sessions.csv", "List of long-term sessions which survive reboots")
	useSSL := flag.Bool("useSSL", false, "If True Enable SSL Server support")
	serverCert := flag.String("serverCert", "server.crt", "Server Cert File")
	serverKey := flag.String("serverKey", "server.key", "Server Key File")

	port := flag.Int("port", 3000, "Port to user")

	flag.Parse()

	// Dump parameters
	log.Printf("\n\tsessionSecret: %v\n\tsessionPeriod: %v\n\tuserFile: %v\n\tsessionFile: %v\n\tuseSSL: %v\n\tserverCert: %v\n\tserverKey: %v\n",
		*services.SessionSecret, *services.SessionPeriod, *services.UserFile, *services.SessionsFile, *useSSL, *serverCert, *serverKey)

	// Load user DB
	if services.LoadUsers(*services.UserFile) != nil {
		log.Println("Error in user loading - abending")
		return
	}
	// Load long term sessions
	if services.LoadSessions(*services.SessionsFile) != nil {
		log.Println("Error in session loading - abending")
		return
	}

	// Start things
	mux := http.NewServeMux()

	// INSTANCE ADMIN SERVICE
	a := rpc.NewServer()
	a.RegisterCodec(json.NewCodec(), "application/json")
	a.RegisterService(new(services.StatsService), "stats")
	a.RegisterService(new(services.HealthService), "health")

	// SESSION SERVICE
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(new(services.SessionService), "session")

	// AUTHENTICATION SERVICE
	u := rpc.NewServer()
	u.RegisterCodec(json.NewCodec(), "application/json")
	u.RegisterService(new(services.AuthenticationService), "authentication")

	// DEFINE ENDPOINTS
	mux.Handle("/api/v1/session", s)
	mux.Handle("/api/v1/admin", a)
	mux.Handle("/api/v1/authentication", u)

	n := negroni.Classic()
	// Stats runs across all instances
	n.Use(services.StatsMiddleware)
	n.UseHandler(mux)

	var err error
	if *useSSL {
		log.Println("Starting in SSL HTTPS Server Mode")
		err = http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), *serverCert, *serverKey, n)

	} else {
		log.Println("Starting in HTTP Server Mode - Passwords can be read by man in the middle.")
		err = http.ListenAndServe(fmt.Sprintf(":%d", *port), n)
	}

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
