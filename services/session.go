package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/muesli/cache2go"
)

type Session struct {
	Id      string   `json:"id"`
	User    string   `json:"user"`
	Expires int64    `json:"expires"`
	Roles   []string `json:"roles"`
}

type SessionInterface interface {
	valid() bool
}

func (s *Session) valid() bool {
	now := time.Now().Unix()
	if s.Expires < now {
		return false
	}
	return true
}

func (s *Session) ToString() string {
	return fmt.Sprintf("SESSION{ ID:%v,\n USER:%v,\n ROLES:%v,\n EXPIRES:%v\n}\n", s.Id, s.User, s.Roles, time.Unix(s.Expires, 0))
}

func init() {
	// cache2go supports a few handy callbacks and loading mechanisms.
	sessions.SetAboutToDeleteItemCallback(func(e *cache2go.CacheItem) {
		created := e.CreatedOn().Format("2006-01-02 15:04:05")
		log.Printf("Session Expired: %v  %v  %v %v left \n", e.Key(), e.Data().(*Session).User, created, (sessions.Count() - 1))
	})
}

// Create a new session for user

type SessionService int

// NewSession - creates a session and add its to the session DB
func (t *SessionService) NewSession(user User) (string, error) {
	log.Printf("NewSession Called for '%v' and secret %v\n", user.User, *SessionSecret)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	exp := time.Now().Add(time.Second * time.Duration(*SessionPeriod)).Unix()
	tokenString, _ := CreateSession(user, *SessionPeriod, *SessionSecret)
	session := Session{tokenString, user.User, exp, user.Roles}
	sessions.Add(tokenString, time.Duration(*SessionPeriod)*time.Second, &session)

	return tokenString, nil

}

// CreateSession - creates a sesion WITHOUT adding it to the session DB
func CreateSession(user User, secondsToLive int, secret string) (string, error) {
	//	log.Printf("CreateSession Called for '%v' and secret %v\n", user.User, secret)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	exp := time.Now().Add(time.Second * time.Duration(secondsToLive)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.User,
		"exp":   exp,
		"jid":   newUUID(),
		"roles": user.Roles,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return tokenString, nil

}

// Look up session
func (t *SessionService) lookup(id string) (Session, error) {
	v, error := sessions.Value(id)
	if error != nil {
		return Session{}, errors.New("204 unknown session")
	}
	if v != nil {
		// OK session known - is it still valid?
		session := v.Data().(*Session)
		if !session.valid() {
			// Expired
			log.Printf("TODO - Remove expired tokens ... in this case %v\n", id)
			return Session{}, errors.New("204 Unknown token")
		} else {
			return *session, nil
		}
	}
	return Session{}, errors.New("204 Unknown token")
}

// Decode
func Decode(tokenString, secret string) (Session, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != "HS256" {
			return Session{}, errors.New("Unsupported Method")
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		return Session{}, err
	}
	//log.Println(token)
	session := Session{}
	session.Id = tokenString
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		session.User = claims["sub"].(string)
		session.Expires = int64(claims["exp"].(float64))
		roles := make([]string, 0)
		croles := claims["roles"].([]interface{})
		for _, v := range croles {
			roles = append(roles, v.(string))
		}
		session.Roles = roles
	}
	return session, nil
}

// Does user has role
func (t *SessionService) hasRole(role string, session Session) bool {

	for _, r := range session.Roles {
		if r == role {
			return true
		}
	}

	return false
}

// Call data
type SessionValidateArgs string
type SessionValidateReply bool

// Verifies is a token is valid - public
func (t *SessionService) Validate(r *http.Request, args *SessionValidateArgs, result *SessionValidateReply) error {
	*result = false

	_, err := t.lookup(string(*args))
	if err == nil {
		// OK session known and valid
		*result = true
	}
	return nil
}

type SessionInvalidateArgs string
type SessionInvalidateReply bool

// invalidates a session
func (t *SessionService) Invalidate(r *http.Request, args *SessionInvalidateArgs, result *SessionInvalidateReply) error {
	*result = false
	var err error

	log.Printf("Invalidating session %s\n", *args)
	if sessions.Exists(string(*args)) {
		_, err = sessions.Delete(string(*args))
		if err == nil {
			// OK session known and valid
			*result = true
			return nil
		}
	} else {
		err = errors.New("204 No Such Session")
	}
	return err
}

// SessionListArgs a struct with a field 'authorization' containing callers token.
// Requires caller to have admin role
type SessionListArgs struct {
	Authorization string `json:"authorization"`
} // This is the callers sessionid
type SessionListReply struct {
	Size int      `json:"size"`
	Ids  []string `json:"ids"`
}

// List available sessions
func (t *SessionService) List(r *http.Request, args *SessionListArgs, result *SessionListReply) error {
	// Check session is valid and has 'admin' role
	val, err := t.lookup(args.Authorization)
	if err != nil {
		// OK they have a problem
		return err
	}
	if !t.hasRole("admin", val) {
		return errors.New("403 You dont have permissions to do that")
	}

	// OK they do have permissions - list
	var sessionKeys []string

	sessions.Foreach(func(key interface{}, item *cache2go.CacheItem) {
		session, _ := item.Data().(*Session)

		if session.valid() {
			sessionKeys = append(sessionKeys, session.Id)
		}
	})
	reply := SessionListReply{sessions.Count(), sessionKeys}
	*result = reply
	return nil
}

// SessionDetailsArgs a struct with a field 'authorization' containing callers token and 'token' which
// is the token to decode. Requires caller to have admin role and
type SessionDetailsArgs struct {
	Authorization string `json:"authorization"`
	Token         string `json:"token"`
}

// List available sessions
func (t *SessionService) Details(r *http.Request, args *SessionDetailsArgs, result *Session) error {
	// Check session is valid and has 'admin' role
	val, err := t.lookup(args.Authorization)
	if err != nil {
		// OK they have a problem
		return err
	}
	if !t.hasRole("admin", val) {
		return errors.New("403 You dont have permissions to do that")
	}
	// They have permission - lookup token
	val, err = t.lookup(args.Token)
	if err != nil {
		// OK they have a problem
		return err
	}
	*result = val
	return nil
}

// Initiaizes data structues - IE Read sessions DB
func LoadSessions(filename string) error {
	csvfile, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer csvfile.Close()
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Reading Sessions Database %s\n", filename)
	// Create user map
	numberOfSessions := 0
	for _, row := range records {
		if len(row) > 0 {
			// Have session so we need to add it to the session DB
			session, err := Decode(row[0], *SessionSecret)
			if err == nil {
				time2LiveSeconds := unixTimeToTTL(session.Expires)
				sessions.Add(session.Id, time.Duration(time2LiveSeconds)*time.Second, &session)
				numberOfSessions = numberOfSessions + 1
				//log.Println(session.ToString())
				log.Printf("\tAdding Session: user[%v] roles%v expires[ %v ] TTL[%v]\n", session.User, session.Roles, time.Unix(session.Expires, 0), time2LiveSeconds)
			} else {
				log.Println("\t[ERROR] ", err)
			}
		}
	}
	log.Printf("#Number of sessions = %v\n", numberOfSessions)
	return nil
}
