BIN=/usr/local/bin
TARGET=gopaste

all: main.go
	ln -s $(shell pwd)/vendor $(shell pwd)/gopath/src 
	GOPATH=$(shell pwd)/gopath go build

install: gopaste
	cp gopaste cleanup.sh $(BIN)/
	setcap 'cap_net_bind_service=+ep' $(BIN)/$(TARGET)

clean:
	go clean

gopaste: all

