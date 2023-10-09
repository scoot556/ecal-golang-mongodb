#!/bin/bash

MONGO_USER="root"
MONGO_PASSWORD="example"
MONGO_HOST="mongo"
MONGO_AUTH_DB="admin"
DATABASE_NAME="ecal"
MONGO_ADMIN_USER="ecalUser"
MONGO_ADMIN_PASSWORD="example"

COLLECTIONS=("comments" "movies")

mongo -u "$MONGO_USER" -p "$MONGO_PASSWORD" --host "$MONGO_HOST" --authenticationDatabase "$MONGO_AUTH_DB" --eval "use $DATABASE_NAME; db.createCollection('comments'); db.createCollection('movies');"

mongo -u "$MONGO_USER" -p "$MONGO_PASSWORD" --host "$MONGO_HOST" --authenticationDatabase "$MONGO_AUTH_DB" --eval "db.getSiblingDB('$DATABASE_NAME').createUser({user: '$MONGO_ADMIN_USER', pwd: '$MONGO_ADMIN_PASSWORD', roles: ['readWrite']});"
echo "User '$MONGO_ADMIN_USER' created with read and write privileges for database '$DATABASE_NAME'."

for collection in "${COLLECTIONS[@]}"; do
  count=$(mongo -u "$MONGO_ADMIN_USER" -p "$MONGO_ADMIN_PASSWORD" --host "$MONGO_HOST" --authenticationDatabase "$DATABASE_NAME" --eval "db.getSiblingDB('$DATABASE_NAME').$collection.count()" --quiet)
  if [ "$count" -gt 0 ]; then
    echo "Collection '$collection' already contains data. Skipping import."
  else
    echo "Collection '$collection' is empty. Importing data..."
    mongoimport --host "$MONGO_HOST" --username "$MONGO_ADMIN_USER" --password "$MONGO_ADMIN_PASSWORD" --authenticationDatabase "$DATABASE_NAME" --db "$DATABASE_NAME" --collection "$collection" --type json --file "./$collection.json"
    echo "Data imported into collection '$collection'."
  fi
done
