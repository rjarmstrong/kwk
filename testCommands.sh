#!/usr/bin/env bash

CONTAINER=localtest_mysql_1

 docker exec -i ${CONTAINER} sh -c 'mysql -uroot -D kwk -e "DELETE FROM aliases"'
 docker exec -i ${CONTAINER} sh -c 'mysql -uroot -D kwk -e "DELETE FROM users"'
 go test -v ./libs/integration/