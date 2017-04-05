VERSION=`cat ./VERSION`

LDFLAGS=-ldflags "-X main.Version=${VERSION}"

install:
	go install -v ${LDFLAGS}

build:
	go build -o ernest -v ${LDFLAGS}

test:
	go test -v ./...

cover:
	go test -coverprofile cover.out

deps:
	go get -u github.com/fatih/color
	go get -u github.com/urfave/cli
	go get -u github.com/mitchellh/go-homedir
	go get -u gopkg.in/yaml.v2
	go get -u github.com/howeyc/gopass
	go get -u github.com/r3labs/sse
	go get -u github.com/olekukonko/tablewriter
	go get github.com/pmezard/go-difflib/difflib
	go get github.com/skratchdot/open-golang/open
	go get github.com/hokaccha/go-prettyjson
	go get github.com/nu7hatch/gouuid

dev-deps: deps
	go get -u github.com/golang/lint/golint
	go get -u github.com/gorilla/mux
	go get -u github.com/smartystreets/goconvey/convey
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/ernestio/ernest-config-client
	go get -u github.com/ernestio/crypto

lint:
	golint ./...
	go vet ./...

dist: dist-linux dist-darwin dist-windows

dist-linux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ernest-${VERSION}-linux-amd64
	zip ernest-${VERSION}-linux-amd64.zip ernest-${VERSION}-linux-amd64 README.md LICENSE
	GOOS=linux GOARCH=386 go build ${LDFLAGS} -o ernest-${VERSION}-linux-386
	zip ernest-${VERSION}-linux-386.zip ernest-${VERSION}-linux-386 README.md LICENSE

dist-darwin:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ernest-${VERSION}-darwin-amd64
	zip ernest-${VERSION}-darwin-amd64.zip ernest-${VERSION}-darwin-amd64 README.md LICENSE
	GOOS=darwin GOARCH=386 go build ${LDFLAGS} -o ernest-${VERSION}-darwin-386
	zip ernest-${VERSION}-darwin-386.zip ernest-${VERSION}-darwin-386 README.md LICENSE

dist-windows:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ernest-${VERSION}-windows-amd64.exe
	zip ernest-${VERSION}-windows-amd64.zip ernest-${VERSION}-windows-amd64.exe README.md LICENSE
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o ernest-${VERSION}-windows-386.exe
	zip ernest-${VERSION}-windows-386.zip ernest-${VERSION}-windows-386.exe README.md LICENSE

clean:
	go clean
	rm -rf ernest-*
