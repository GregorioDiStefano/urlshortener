#!/bin/bash
# Start backend

set -e

docker-compose down --remove-orphans && docker-compose up -d
sleep 10 # the app will restart a few times until the db and redis are up, there are better ways of doing this.

app=$(docker ps --format '{{.Names}}' | grep 'urlshort' | grep 'app')
backend_ip=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $app)

echo "Backend is running at http://$backend_ip:8888 and on http://localhost:8888"
