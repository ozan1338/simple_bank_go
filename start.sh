#!/bin/sh

set -e

echo "run db migration"
# echo "$USER"
# whoami
# ls -al /app
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"