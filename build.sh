export GOPATH=`pwd`
os=$1
arch=$2
GOOS=$os GOARCH=$arch go build -o bin/sms${os}${arch} -ldflags \
"-X main.buildTime=`date +%Y-%m-%d/%H:%M:%S`" src/main/*.go

