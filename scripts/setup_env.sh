#!/bin/bash

# Check if environment argument is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <environment>"
    echo "Available environments: local, dev, prod"
    exit 1
fi

ENV=$1
ENV_FILE="config/env.$ENV"

# Check if environment file exists
if [ ! -f "$ENV_FILE" ]; then
    echo "Environment file $ENV_FILE not found"
    exit 1
fi

# Create secrets directory if it doesn't exist
mkdir -p secrets/$ENV

# Load environment variables
export $(cat $ENV_FILE | grep -v '^#' | xargs)

# Create secret files if they don't exist (for dev and prod environments)
if [ "$ENV" != "local" ]; then
    # Database secrets
    echo "Creating secret files for $ENV environment..."
    
    if [ ! -f "secrets/$ENV/db_host" ]; then
        read -p "Enter database host for $ENV: " db_host
        echo "$db_host" > "secrets/$ENV/db_host"
    fi

    if [ ! -f "secrets/$ENV/db_user" ]; then
        read -p "Enter database user for $ENV: " db_user
        echo "$db_user" > "secrets/$ENV/db_user"
    fi

    if [ ! -f "secrets/$ENV/db_password" ]; then
        read -s -p "Enter database password for $ENV: " db_password
        echo
        echo "$db_password" > "secrets/$ENV/db_password"
    fi
fi

echo "Environment $ENV has been set up successfully" 