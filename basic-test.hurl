POST http://localhost:8888/auth/register
{
    "email": "greg.distefano@gmail.com",
    "password": "blahblahblah!234"
}
HTTP 201


POST http://localhost:8888/auth/login
{
    "email": "greg.distefano@gmail.com",
    "password": "blahblahblah!234"
}
HTTP 200
[Captures]
token: jsonpath "$['token']"

POST http://localhost:8888/api/v1/shorten
Authorization: Bearer {{token}}
{
    "url": "https://www.google.com"
}
HTTP 201
[Captures]
shortUrl: jsonpath "$['key']"

GET http://localhost:8888/{{shortUrl}}
HTTP 301
Location: https://www.google.com

GET http://localhost:8888/api/v1/urls/
Authorization: Bearer {{token}}
HTTP 200