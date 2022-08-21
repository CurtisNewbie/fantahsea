#!/bin/bash

# --------- remote
remote="alphaboi@curtisnewbie.com"
remote_path="~/services/fantahsea/build"
# ---------

#scp "run.sh" "${remote}:${remote_path}"
#scp "Dockerfile" "${remote}:${remote_path}"

# copy the config file just in case we updated it
#scp "app-conf-dev.json" "${remote}:${remote_path}"

GLOBIGNORE='fantahsea-tmp'

scp -r * "${remote}:${remote_path}/"
