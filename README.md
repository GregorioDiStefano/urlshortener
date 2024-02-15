# Basic Short URL

### A basic URL shorter backend service

Start backend (NOTE: you need `docker` and `docker-compose` installed):

`./start-backend.sh`

#### Usage:

`curl -v -X POST http://localhost:8888/auth/register -H "Content-Type: application/json" -d '{"email": "user@example.com", "password": "your_password"}'`

`curl -v -X POST http://localhost:8888/auth/login -H "Content-Type: application/json" -d '{"email": "user@example.com", "password": "your_password"}'`

`curl -v -X POST http://localhost:8888/api/v1/shorten -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" -d '{"url": "https://example.com"}'`

`curl -v -X GET http://localhost:8888/api/v1/urls/ -H "Authorization: Bearer <TOKEN>"`

`curl -v -X DELETE http://localhost:8888/api/v1/shorten/YOUR_SHORT_URL_KEY -H "Authorization: Bearer <TOKEN>"`

#### Design considerations:

* Written in Go
* Heavily unit/integration tested with 72% coverage (used interfaces to allow mocking, but relied on using a test db instead since it add more value)
* E2E testing using hURL (https://hurl.dev/)
* PostgreSQL for persistent data storage + Redis for caching
* JWT based authenticated
* Added a health endpoint to check postgresql and redis are happy (useful for k8s, docker, other health checkers)
* Instead of using a traditional base64 on database row ID, I wanted to make enumeration/url guessing slightly more difficult by suffixing the short url key with a random value (which I incorrectly called a 'nonce'..). That means your should url is basically base64_without_padding(rowid) + 2 random characters  
* I simplified a lot of things (config, db setup), wrote tests to test basic happy paths and common errors to save time. 
* Passwords like "mysecretpassword" / "hunter2" are not intended to be serious passwords.
* OpenAPI spec provided.

#### Left to do (just to name a few):

* TODOs
* Better logging / log statements / error messages
* Improve unit tests
* Add instrumentation / observability (prometheus?)
* Migration support
* Admin endpoint to administrator users
* Setup k8s config
* Get rid of utils.go
* etc etc etc
