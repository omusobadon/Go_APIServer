@echo off
SETLOCAL

echo Go_APIServer downloading...
git clone https://github.com/omusobadon/Go_APIServer.git

echo Changeing directory...
cd Go_APIServer

git pull
git switch local-playing
git pull

timeout /t 3

echo service_app is createing..
podman compose -f .\docker-compose.yml up -d

echo  wait for starting service_app ...
timeout /t 3


echo Go_APIServer Setuping...
go get github.com/steebchen/prisma-client-go
go run github.com/steebchen/prisma-client-go db push

echo seting...
timeout /t 3

echo Go_APIServer is started!
go run .

#echo service_app is started!
#echo Go_APIServer is started!
#podman build -t api-server .
#
#echo wait for starting Go_APIServer ...
#timeout /t 3
#podman run -d -p 8080:8080 api-server
#
#echo wait for starting Go_APIServer ...
#timeout /t 3

