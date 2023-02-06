BINARY=stream-go
all: build

build:
	go build -o $(BINARY)

clean:
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi

run: build
	./$(BINARY)

