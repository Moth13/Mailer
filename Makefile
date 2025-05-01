templ:
	templ generate

test:
	go test -v -cover ./...

air:
	air

server:
	go run main.go
