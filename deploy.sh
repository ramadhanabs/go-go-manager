#!/bin/bash

# Load environment variables from .env file
set -o allexport
source .env
set +o allexport

# Remove existing compiled file
rm -rf $APP_NAME.gz

# Build the Go binary
GOOS=linux GOARCH=amd64 go build -o $APP_NAME

# Compress the binary
gzip $APP_NAME

FILE_TO_TRANSFER="${APP_NAME}.gz"
SERVERS=("$SERVER_INSTANCE_1" "$SERVER_INSTANCE_2")

# Loop through each server
for SERVER in "${SERVERS[@]}"; do
    echo "Deploying to $SERVER..."

    # Clean up existing files on the server before deploying
    ssh $SERVER -i $SSH_KEY_PATH << EOF
        cd $REMOTE_DIR
        echo "Cleaning up old build files..."
        rm -f $APP_NAME
        rm -f $FILE_TO_TRANSFER
EOF

    # SCP the new file to the server
    echo "Transferring new build file to $SERVER..."
    scp -i $SSH_KEY_PATH $FILE_TO_TRANSFER $SERVER:$REMOTE_DIR/

    # SSH into the server to decompress and set up the application
    ssh $SERVER -i $SSH_KEY_PATH << EOF
        cd $REMOTE_DIR
        echo "Setting up the new build..."
        gunzip $FILE_TO_TRANSFER
        chmod +x $APP_NAME
        echo "Deployment completed on $SERVER!"
EOF

    echo "Deployment completed for $SERVER."
done