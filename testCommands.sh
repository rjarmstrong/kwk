#!/usr/bin/env bash

 docker exec -i mysql_test sh -c 'mysql -uroot -D kwk -e "DELETE FROM aliases"'
 docker exec -i mysql_test sh -c 'mysql -uroot -D kwk -e "DELETE FROM users"'
 go test ./libs/integration/