#!/bin/bash

# --------- remote
remote="alphaboi@curtisnewbie.com"
remote_path="~/services/fantahsea/build"
# ---------

#scp "run.sh" "${remote}:${remote_path}"
#scp "Dockerfile" "${remote}:${remote_path}"

# copy the config file just in case we updated it
#scp "app-conf-dev.json" "${remote}:${remote_path}"

GLOBIGNORE="fantahsea-tmp:.git:.vscode:fantahsea-base"

ssh  "alphaboi@curtisnewbie.com" "mv services/fantahsea/build services/fantahsea/build-bak"

scp app-conf-prod.yml "${remote}:~/services/fantahsea/config"
scp -r * "${remote}:${remote_path}/"
if [ ! $? -eq 0 ]; then
    exit -1
fi

ssh  "alphaboi@curtisnewbie.com" "cd services; docker-compose up -d --build fantahsea"
