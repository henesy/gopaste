BIN=/usr/local/bin
TARGET=gopaste

all: main.go
	go build

install: gopaste
	cp gopaste cleanup.sh $(BIN)/

clean:
	go clean

gopaste: all

