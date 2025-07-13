b:
	go build -o build/glox

r:
	./build/glox

t:
	go build -o build/glox
	go test ./test/...

dev:
	go build -o build/glox
	./build/glox
