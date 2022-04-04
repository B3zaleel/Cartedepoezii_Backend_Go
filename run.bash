#!/usr/bin/env bash
declare -A ENV_VARS
File_Lines=()

# read environment variables from file (.env)
readarray -t File_Lines < <(cat .env)
for ((i = 0; i < "${#File_Lines[@]}"; i++)) do
    line="${File_Lines[i]}"
    ENV_VARS["$(echo "$line" | cut -d ':' -f1)"]="$(echo "$line" | cut -d ' ' -f2-)"
done

env GOPATH="$PWD/src" \
    GIN_MODE="${ENV_VARS['GIN_MODE']}" \
    DB_URL="${ENV_VARS['DB_URL']}" \
    APP_MAX_SIGNIN_TRIES="${ENV_VARS['APP_MAX_SIGNIN_TRIES']}" \
    HOST="${ENV_VARS['HOST']}" \
    IMG_CDN_PUB_KEY="${ENV_VARS['IMG_CDN_PUB_KEY']}" \
    IMG_CDN_PRI_KEY="${ENV_VARS['IMG_CDN_PRI_KEY']}" \
    IMG_CDN_URL_EPT="${ENV_VARS['IMG_CDN_URL_EPT']}" \
    GOOGLE_MAIL_SENDER="${ENV_VARS['GOOGLE_MAIL_SENDER']}" \
    WEB_CLIENT_DOMAIN="${ENV_VARS['WEB_CLIENT_DOMAIN']}" \
    PWD_SALT="${ENV_VARS['PWD_SALT']}" \
    APP_SECRET_KEY="${ENV_VARS['APP_SECRET_KEY']}" \
    go run src/main.go
