@echo off
SETLOCAL

echo PostgreSQL is createing..
podman compose -f .\docker-compose.yml up -d service_postgres

echo PostgreSQL is wait for starting...
timeout /t 3

echo pgAdmin is createing
podman compose -f .\docker-compose.yml up -d service_pgadmin

echo pgAdmin is wait for starting...
timeout /t 3

echo service_app is createing..
podman compose -f .\docker-compose.yml up -d service_app

echo  wait for starting service_app ...
timeout /t 3

echo Done!
ENDLOCAL
