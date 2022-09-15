#!/bin/sh

set -e

echo "run db migration"
# echo "$USER"
# whoami
# ls -al /app
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"