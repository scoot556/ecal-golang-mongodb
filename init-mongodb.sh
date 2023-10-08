#!/bin/bash
set -e

# Wait for MongoDB to start
until mongo --eval "quit(db.runCommand({ ping: 1 }).ok ? 0 : 1)"; do
  echo "Waiting for MongoDB to start..."
  sleep 1
done

# Create the "ecal" database
mongo admin --eval "db.getSiblingDB('ecal')"

# Import data from "comments.json" into "comments" collection
mongoimport --host mongodb --db ecal --collection comments --type json --file /data/comments.json --jsonArray

# Import data from "movies.json" into "movies" collection
mongoimport --host mongodb --db ecal --collection movies --type json --file /data/movies.json --jsonArray
