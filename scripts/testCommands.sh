#!/usr/bin/env bash

docker exec -i cass cqlsh -e "use kwk; TRUNCATE snips; TRUNCATE users_by_email; TRUNCATE users;"
go test -short ./libs/integration/

-- snips, users_by_email, users