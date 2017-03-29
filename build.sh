export GOPATH=`pwd`
os=$1

GOOS=$os go build -o bin/sms$os -ldflags \
"-X main.buildTime=`date +%Y-%m-%d/%H:%M:%S`" src/main/*.go

