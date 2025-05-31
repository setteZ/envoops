#! /bin/bash
if [ -z "$1" ]; then
    echo "Error: Missing argument." >&2
    exit 1
fi

if [ -z "$SSH_USER" ]; then
    echo "Error: Missing \$SSH_USER." >&2
    exit 1
fi

if [ -z "$SERVER_IP_NODE_OTA" ]; then
    echo "Error: Missing \$SERVER_IP_NODE_OTA." >&2
    exit 1
fi

if [ -z "$NODE_RELEASES_PATH" ]; then
    echo "Error: Missing \$NODE_RELEASES_PATH." >&2
    exit 1
fi

if [ ! -d "release-node/$1" ]; then
    echo "$1 does not exist in folder release-node"
    exit 1
fi

if [ -n "$SSH_PRIVATE_KEY_LOCATION" ]; then
    ID_FILE="-i $SSH_PRIVATE_KEY_LOCATION"
else
    ID_FILE=""
fi

scp $ID_FILE -r release-node/* $SSH_USER@$SERVER_IP_NODE_OTA:$NODE_RELEASES_PATH
ssh $ID_FILE $SSH_USER@$SERVER_IP_NODE_OTA "cd $NODE_RELEASES_PATH && echo $1 > version"
