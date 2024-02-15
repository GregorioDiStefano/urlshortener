#!/bin/bash
# Start backend

set -e

docker-compose down --rmi all --remove-orphans && docker-compose up -d --force-recreate 
sleep 10 # the app will restart a few times until the db and redis are up, there are better ways of doing this.

backend_ip=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' urlshort-app-1)

echo "Backend is running at http://$backend_ip:8888 and on http://localhost:8888"
