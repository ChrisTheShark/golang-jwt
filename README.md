# go-jwts
Simple Go implementation of basic authentication providing a JWT and authorization using a parsed JWT. My purpose for this project was to POC JWT functionality using middleware in selective chi routes. Not all routes in the proposed finalized project will be secured using JWTs. Currently only passing down the user name using the request context but roles could be passed as well to authorize.  

## Getting Started
First, clone this project or execute `go get` to retrieve project files. Navigate to directory and execute `go run main.go`.

## Authenticating for a JWT
Currently, to simplify the implementation there is a single static user and password. Feel free to update these values in main.go to reflect a user of your choice. The password is hashed upon startup and the hash is compared to the incoming value to authenticate. This is not meant to be a true battle tested login mechanism, we need to simplify the implementation to ease comprehension. 


```
curl -X POST \
  http://localhost:8080/token \
  -H 'content-type: application/json' \
  -d '{
	"username": "chrisd",
	"password": "hotdog23"
}'
```

The curl command above should result in a JWT response similar to the example below.

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTI3MDk3NjA4LCJuYW1lIjoiQ2hyaXMgRHllciJ9.9Db-ny27fh2fo_OQB2AWL1KlxpDK6qXqbAdcYuvPhA
```

## Authorization with a JWT
Currently there is a single route which requires a JWT `/greeting`. Provide the recieved JWT as an Authorization header to retrieve a greeting specific to the authenticated user.

Example `/greeing` request:

```
curl -X GET \
  http://localhost:8080/greeting \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTI3MDkwNTM5LCJuYW1lIjoiQ2hyaXMgRHllciJ9.flbTbmlqJalHLmlfkUrWrykmc5q8S990XRZGGCQozWk' \
```
