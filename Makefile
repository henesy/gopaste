BIN=/usr/local/bin
TARGET=gopaste

all: main.go
	GOPATH=$(shell pwd)/gopath go build -mod=vendor

install: gopaste
	cp gopaste cleanup.sh $(BIN)/
	setcap 'cap_net_bind_service=+ep' $(BIN)/$(TARGET)

clean:
	go clean

gopaste: all

# For paste.iseage.org running go1.10.x/amd64
setup:
	mkdir ./gopath
	ln -s $(shell pwd)/vendor $(shell pwd)/gopath/src 

