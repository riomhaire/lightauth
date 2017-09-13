# lightauth

This project *SHOULD NOT* be used in production for various reasons. This project is one of our many experiments in learning GO and in this case JSON-RPC, Crypto, Gorilla RPC and JWT for a demo Authentication and Session application. We prefer to learn using something more than the usual 'hello world'.

In essence the API allows a UI or other API to authenticate a user with username/password (passwords are hashed so no leak) and returns a 'token' which can be used and submitted when it can be queried for validity or checked for claims (if 'admin' token is also supplied).

All communication is done via json RPC.

## Installation

The simplest way if you 'make' installed is to run 'make' which will install all the dependencies and install the apps. Otherwise After cloning the repository you need to install the few dependencies. Execute the following within the main directory.

```bash
$ go get github.com/riomhaire/lightauth
$ cd <gopath-root>/src/github.com/riomhaire/lightauth
$ go get ./...
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
$ lightauth --help
lightauth version 0.4
Usage of lightauth:
  -port int
    	Port to user (default 3000)
  -serverCert string
    	Server Cert File (default "server.crt")
  -serverKey string
    	Server Key File (default "server.key")
  -sessionFile string
    	List of long-term sessions which survive reboots (default "sessions.csv")
  -sessionPeriod int
    	How many seconds before sessions expires (default 3600)
  -sessionSecret string
    	Master key which is used to generate system jwt (default "secret")
  -useSSL
    	If True Enable SSL Server support
  -usersFile string
    	List of Users and salted/hashed password with their roles (default "users.csv")
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

The API is a fairly simple one and consists of:

1. Authenticate/Login.
2. Verify a session token.
3. Invalidate session.
4. List known session tokens.
5. Get session token details (get roles).
6. Get Call Statistics.

The endpoint in the default startup can be found at "http://somehost:3000/api/v1/authentication" or "http://somehost:3000/api/v1/session" or "http://somehost:3000/api/v1/admin" - and yes we know we should be using HTTPS or some other transport medium, but this a simple project to help us learn GO and not a prod app.

Content-Type should be "application/json"

### Authenticate/Login.

Example request ("http://somehost:3000/api/v1/authentication"):

```json
{
	 "id":"-1",
	  "method":"authentication.Authenticate",
	  "params":[{
			"username":"vader",
			"password":"anakin"
		}]
}
```
Wont mention password is in clear - should be a hash or some other method in conjunction with https... but this is a learning experience app.

Example response on success:

```json
{
	"result": {
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDQ5Nzk4OTUsImppZCI6IjVmODM1ZGYxLTkaZWYtNGRisC1hNTRlLTMyZWFiMThkNDJhMCIsInJvbGVzIjpbInNlY3VyaXR5Il0sInN1YiI6InN0b3JtdHJvb3BlciJ9.XtBeitAsJB63Xj_6q6BqWU0Lo0qgBhHz7oSStPZxchI"
	},
	"error": null,
	"id": "-1"
}

```

Example response on error:
```json
{
	"result": null,
	"error": "401 Not Allowed",
	"id": "-1"
}

```

### Verify a session token.

Example request ("http://somehost:3000/api/v1/session"):

```json
{
	 "id":"-1",
	  "method":"session.Validate",
	  "params":[
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDQ5N4k4OTUsImppZCI6IjVmODM1ZGYxLTk5ZWYtNGRiZC1hNTRlLTMyZWFiMThkNDJhMCIsInJvbGVzIj1bInNlY3VyaXR5Il0sInN1YiI6InN0b3JtdHJvb3BlciJ9.XtBeitAsJB63Xj_6q6BqWU0Lo0qgBhHz7oSStPZxchI"
		]
}

```

Example response:

```json
{
	"result": true,
	"error": null,
	"id": "-1"
}

```
result is true or false depending is session token is valid or not.


### Invalidate session.

Example request  ("http://somehost:3000/api/v1/session"):

```json
{
	 "id":"-1",
	  "method":"session.Invalidate",
	  "params":[
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDQ5Nzk4OTUsImppZCI6IjVmODM1ZGYxLTk5ZWYtNGRiZC1hNTRlLTMyZWFiMThkNDJhMCIsInJvbGVzIjpaInNlY3VyaXR5Il0sInN1YiI6InN0b3JtdHJvb3BlciJ9.XtBeitAsJB63Xj_6q6BqWU0Lo0qgBhHz7oSStPZxchI"
		]
}
```

Example response:

```json
{
	"result": true,
	"error": null,
	"id": "-1"
}

```
### List known session tokens.

Example request ("http://somehost:3000/api/v1/session"):

```json
{
	 "id":"-1",
	  "method":"session.List",
	  "params": [ { "authorization":
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI0MDQ2MzQzNzEsImppqCI6Ijg2MTkwNzVmLTA3ZGUtNDk5Yy1iMTgyLWIyNTJiZDNhYjM3YiIsInJvbGVzIjpbImFkbWluIl0sInN1YiI6ImdyZW1saW4ifQ.oAuj5dw4sg0F8aETC9t8d_LrJe_PXj601SDq4xD6Fig"
								}]
}
```

Note: the 'authorization' token should be one which contains a 'admin' role.


Example response:

```json
{
	"result": {
		"size": 5,
		"ids": [
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI1MDQ5NTM5NDAsImppZCI6ImU3YmMxYWY2LTc0MDMtNDFkNi05NWI2LTUxMDg5ODU1NjZkZCIsInJvbGVzIjpbImFwaSIsImFkbWluIl0sInN1YiI6ImdyZW1saW4ifQ.xKQIxtNXlFxnwV5kE2Z0KWB_hV0cgPGHiNnsJgji494",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDQ5ODQ1MTMsImppZCI6Ijc1ZmM2MWZjLTk4NjItNDZdZC05MTk1LWJmZmY3ZGVlOTFiYSIsInJvbGVzIjpbInNlY3VyaXR5Il0sInN1YiI6InN0b3JtdHJvb3BlciJ9.UnBlzvx4clDMNI5i8mLBub56IuS6hb_lPwpR0bqRSKg",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDQ5ODQ1MTAsImppZCI6ImZhYTA5M2RlLWU0MGE4NDQ0ZC04N2RiLWFhOWU4N2VjMDIyYSIsInJvbGVzIjpbInNlY3VyaXR5Il0sInN1YiI6InN0b3JtdHJvb3BlciJ9.LFajdn1Nx1_fsTsR96sIuswvqHHQBvxcZZbBKRyjW7I",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI0MDQ2MzQzNzEsImppZCI6Ijg2MTkwNzVmLTA1ZGUtNDk5Yy1iMTgyLWIyNTJiZDNhYjM3YiIsInJvbGVzIjpbImFkbWluIl0sInN1YiI6ImdyZW1saW4ifQ.oAuj5dw4sg0F8aETC9t8d_LrJe_PXj601SDq4xD6Fig",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQ5aDM0OTMsImppZCI6IjVkOWIwMmNlLTcyMzItNDQyZS1hNGRiLTkxMzRlMDFiODZhYSIsInJvbGVzIjpbImFwaSJdLCJzdWIiOiJzeXN0ZW0ifQ.a_VrhyXMMBoK0lZp5CnmWz__SiWOVf1KMWUF_ljCl9w"
		]
	},
	"error": null,
	"id": "-1"
}

```
### Get session token details (get roles).

Example request ("http://somehost:3000/api/v1/session"):

```json
{
	 "id":"-1",
	  "method":"session.Details",
	  "params": [ { "authorization":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI0MDQ2MzQzNzEsImppZCI6Ijg2MTkwNzVmLTA3ZGUtNDk5Yy1iMTgyLWIyNTJiZDNhYjM3YiIsInJvbGVzIjpbImFkbWluIl0aInN1YiI6ImdyZW1saW4ifQ.oAuj5dw4sg0F8aETC9t8d_LrJe_PXj601SDq4xD6Fig",
        "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDQ5NzkyMDYsImppZCI6IjNmYWM5MjU2LTNjNTEtNGM5OC05YzZlLWU1MjA1NGMzYzIyZSIsInJvbGVzIjpbInNlY3VyaXR5Il0sInN1YiI6InN0b3JtdHJvb3BlciJ9.4bQh09BPOjaavzhvErIb008Ot6STyd2B-ZoXta-7g2Y"
    }]
}

```
Note: This API requires the caller 'authorization' token to have an encoded 'admin' role.

Example response is successful:

```json
{
	"result": {
		"id": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDQ5NzkyMDYsImppZCI6IjNmYWM5MjU2LTNjNTEtNGM5OC05YzZlLWU1MjA1NGMzYzIyZSIsInJvbGVzIjpbInNlY3VyaXR5Il0sInN1YiI6InN0b3JtdHJvb3BlciJ9.4bQh09BPOjaavzhvErIb008Ot6STyd2B-ZoXta-7g2Y",
		"user": "stormtrooper",
		"expires": 1504979206,
		"roles": [
			"security"
		]
	},
	"error": null,
	"id": "-1"
}
```
The 'error' field will be non-null on error.

### Get Call Statistics.

More of an devops call which we added to see what "github.com/thoas/stats" would produce.

Example request ("http://somehost:3000/api/v1/admin"):

```json
{
	 "id":"-1",
	  "method":"stats.Status",
	  "params":[]
}
```

Example response:

```json
{
	"result": {
		"pid": 3652,
		"uptime": "44m34.810671702s",
		"uptime_sec": 2674.810671702,
		"time": "2017-09-09 19:15:22.892921156 +0100 IST m=+2674.812216312",
		"unixtime": 1504980922,
		"status_code_count": {},
		"total_status_code_count": {
			"200": 142813
		},
		"count": 0,
		"total_count": 142813,
		"total_response_time": "17m53.138264488s",
		"total_response_time_sec": 1073.138264488,
		"average_response_time": "7.514289ms",
		"average_response_time_sec": 0.007514289
	},
	"error": null,
	"id": "-1"
}
```

## SSL

As of 0.4 Support for SSL server has been added based on information at that excellent resource https://gist.github.com/6174/9ff5063a43f0edd82c8186e417aae1dc and is enabled via three command line variables:

* useSSL - set that to 'true' EG '-useSSL true'
* serverCert - contain the name of the file containing the SSL cert to use
* serverKey - contain the name of the file containing the server key.

For self signed certificates you can use the following steps to generate them:

```bash
# Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048
    
# List ECDSA the supported curves (openssl ecparam -list_curves)
openssl ecparam -genkey -name secp384r1 -out server.key

openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```


## Examples

To-Do
