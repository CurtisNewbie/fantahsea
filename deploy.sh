#!/bin/bash

GLOBIGNORE="fantahsea-tmp:.git:.vscode:fantahsea-base"

# --------- remote
remote="alphaboi@curtisnewbie.com"
remote_path="~/services/fantahsea/build"
# ---------


ssh  "alphaboi@curtisnewbie.com" "rm -rv ${remote_path}/*"

scp app-conf-prod.yml "${remote}:~/services/fantahsea/config"
if [ ! $? -eq 0 ]; then
    exit -1
fi

scp -r * "${remote}:${remote_path}/"
if [ ! $? -eq 0 ]; then
    exit -1
fi

ssh  "alphaboi@curtisnewbie.com" "cd services; docker-compose up -d --build fantahsea"
