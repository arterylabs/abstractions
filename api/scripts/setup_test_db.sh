#!/bin/zsh

# Default values for database configuration if not set in environment
: ${DATABASE_HOST:="localhost"}
: ${DATABASE_USER:="postgres"}
: ${DATABASE_PASSWORD:="postgres"}
: ${DATABASE_PORT:=5432}
: ${DATABASE_SSL_MODE:="disable"}
: ${DATABASE_NAME:="abstractions"}
TEST_DB_NAME="${DATABASE_NAME}_test"

# # Load environment variables if .env file exists (these can override the defaults)
# if [ -f .env ]; then
#     source .env
# fi

echo "Setting up test database..."

# Try to connect to PostgreSQL and create test database
PGPASSWORD=$DATABASE_PASSWORD psql -h $DATABASE_HOST -p $DATABASE_PORT -U $DATABASE_USER postgres <<EOF
DROP DATABASE IF EXISTS $TEST_DB_NAME WITH (FORCE);
CREATE DATABASE $TEST_DB_NAME;
CREATE EXTENSION IF NOT EXISTS vector;
EOF

if [ $? -eq 0 ]; then
    echo "Successfully created test database: $TEST_DB_NAME"
    
    echo "Running migrations..."
    go run cmd/main.go migrate \
        --host localhost \
        --port 5432 \
        --user postgres \
        --password postgres \
        --dbname "$TEST_DB_NAME" \
        --sslmode disable \
        0
    
    if [ $? -eq 0 ]; then
        echo "Successfully ran migrations"
        exit 0
    else
        echo "Failed to run migrations"
        exit 1
    fi
else
    echo "Failed to create test database"
    exit 1
fi
