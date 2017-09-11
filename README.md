# lightauth

This project *SHOULD NOT* be used in production for various reasons. This project is one of our many experiments in learning GO and in this case JSON-RPC, Crypto and JWT for a demo Authentication and Session application. We prefer to learn using something more than the usual 'hello world'.

In essence the API allows a UI or other API to authenticate a user with username/password (passwords are hashed so no leak) and returns a 'token' which can be used and submitted when it can be queried for validity or checked for claims (if 'admin' token is also supplied).

All communication is done via json RPC.

## Installation

After cloning the repository you need to install the few dependencies. Execute the following within the main directory.

```bash
go get ./...
```

There are three applications - the server is in the 'lightauth' directory, a session generation application (useful for generating API tokens with roles such as admin), and a 'user' app for creating users suitable for including in the users.csv file (similar to passwd). 

The best way if you just want run is to build and install the apps:

```bash
 go install github.com/riomhaire/lightauth
 go install github.com/riomhaire/lightauth/lightauthsession
 go install github.com/riomhaire/lightauth/lightauthuser

```

## Getting Started


### User Creation Application

The 'users.csv' file is read from the folder where the 'lightauth' server is executed. This is a simple CSV file where the 1st line is a header consisting of:

```csv
username,password,enabled,roles

```
An example is:

```csv
username,password,enabled,roles
test,939c1f673b7f5f5c991b4d4160642e72880e783ba4d7b04da260392f855214a6,true,none
admin,50b911deac5df04e0a79ef18b04b29b245b8f576dcb7e5cca5937eb2083438ba,true,admin

```

The password is hashed based on a secret and a salt.  To add a user you need to use the 'lightauthuser' application which takes parameters and creates a line suitable to append to the user csv file:

```bash
$ lightauthuser --help
Usage of lightauthuser:
  -password string
        Password to use - cannot be empty
  -roles string
        List of roles separated by ':' (default "guest:public")
  -user string
        Username associated with the token (default "anonymous")
```

Roles can be any sensible string 'user', 'api' etc the only one used by the session api is 'admin' which is required for some calls such as 'list known sessions' and 'decode token'.


### Session Creation Application

The lightauth server can read a list of sessions from a file called 'sessions.csv' which is used to store a list of long lived session tokens such as those for API's. If a session token is created via the authenticate method they have a limited life span (usually 3600 seconds) before they become invalid. Tokens for api's typically have a requirement for a longer lived period - sometimes months or longer. Long lived Tokens need not be stored within the sessions file since they only need to be encoded using the same parameters as used by the lightauth server itself. The sessions flle is only really used to query for a list of KNOWN sessions tokens.

The session token creation application is called 'lightauthsession':

```bash
$ lightauthsession --help
Usage of lightauthsession:
  -roles string
        List of roles separated by ':' (default "guest:public")
  -secret string
        Key used to generate sessions (default "secret")
  -sessionPeriod int
        How many seconds before sessions expires (default 3600)
  -token string
        If populated means decode token
  -user string
        Username associated with the token (default "anonymous")
```

An example usage would be:

```bash
 $ lightauthsession -user someapp -roles "api:admin" -sessionPeriod 9999999 -secret hush

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTUxNTYzMTIsImppZCI6IjJiM2ZhYjcwLTA3M2MtNGRiNi05ZTEwLThlOWJlMTQwZWM5NCIsInJvbGVzIjpbImFwaSIsImFkbWluIl0sInN1YiI6InNvbWVhcHAifQ.dmsHOMzspru-HBL7QsuLILhFuEOlNSXMksVUismFi8U
```

Created sessions are only valid for lightauth servers which have been started with the same secret.


### The LightAuth Server

The lightauth is a single server which includes authentication and session token management. In a production system these should be implemented as separate servers - but this is a learning experience and we dont plan to do this at the moment.

The server can be started with the following parameters:

```bash
$ lightauth --help
Usage of lightauth:
  -port int
        Port to user (default 3000)
  -sessionFile string
        List of long-term sessions which survive reboots (default "sessions.csv")
  -sessionPeriod int
        How many seconds before sessions expires (default 3600)
  -sessionSecret string
        Master key which is used to generate system jwt (default "secret")
  -usersFile string
        List of Users and salted/hashed password with their roles (default "users.csv")
```
The parameters are pretty much self evident. An example startup would produce:

```bash
$ lightauth
2017/09/11 20:04:59 Reading User Database users.csv
2017/09/11 20:04:59     User test, Enabled = true
2017/09/11 20:04:59     User admin, Enabled = true
2017/09/11 20:04:59 #Number of users = 2
2017/09/11 20:04:59 Reading Sessions Database sessions.csv
2017/09/11 20:04:59     Adding Session: user[test] roles[none] expires[ 2020-11-12 03:34:26 +0000 GMT ] TTL[99995366]
2017/09/11 20:04:59 #Number of sessions = 1
[negroni] listening on :3000
```

## The API

