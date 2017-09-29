CURRENT=pwd
cd /usr/local/src/slack-lambda
go build -o bin/proxy-server proxy-server/main.go 
go build -o bin/container-server container-server/main.go 
cd $CURRENT
