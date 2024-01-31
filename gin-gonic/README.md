A simple example of REST API written in the Golang programming language. To implement routing, the Gin-gonic framework was used. Resources are stored in a file inside the /data folder.

Example:

POST:
curl --location 'http://localhost:8080/players' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email" : "player1@example.com"
}'

GET:
curl --location 'http://localhost:8080/players/{playerId}'

curl --location 'http://localhost:8080/players'

DELETE:
curl --location --request DELETE 'http://localhost:8080/players/{playerId}'
