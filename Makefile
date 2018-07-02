GIT_VERSION=`git log --pretty=format:"%h" -1`
BIN_VERSION=`cat version.txt`

osx:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.gitCommit=${GIT_VERSION} -X main.appVersion=${BIN_VERSION}" -o proxy-ng-darwin main.go

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.gitCommit=${GIT_VERSION} -X main.appVersion=${BIN_VERSION}" -o proxy-ng-linux main.go

windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.gitCommit=${GIT_VERSION} -X main.appVersion=${BIN_VERSION}" -o proxy-ng-windows main.go

build: osx linux windows