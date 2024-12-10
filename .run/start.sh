#!/use/bin/bash

set -a
. .env
set +a

go run ../cmd/cleaner/main.go