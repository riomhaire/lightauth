package services

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

// Create a new session for user

type SessionService int

// NewSession - creates a session and add its to the session DB
func (t *SessionService) NewSession(user User) (string, error) {
	log.Printf("NewSession Called for '%v' and secret %v\n", user.User, *SessionSecret)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	tokenString, _ := CreateSession(user, *SessionPeriod, *SessionSecret)
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
func (session *Session) hasRole(role string) bool {

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
	session, err := Decode(string(*args), *SessionSecret)
	if err != nil {
		// OK they have a problem - so say no. Error not important
		return nil
	}
	// If session not expired
	if session.valid() {
		*result = true
	}
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
	// Check caller session is valid and has 'admin' role
	session, err := Decode(args.Authorization, *SessionSecret)
	if err != nil {
		// OK they have a problem
		return err
	}
	if !session.hasRole("admin") || session.hasRole("api") {
		return errors.New("403 You dont have permissions to do that")
	}
	// They have permission - decode token
	session, err = Decode(args.Authorization, *SessionSecret)
	if err != nil {
		// OK they have a problem
		return err
	}
	*result = session
	return nil
}

// Get version
func (t *SessionService) Version(r *http.Request, args *string, result *string) error {
	*result = Version
	return nil
}
