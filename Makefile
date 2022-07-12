
restdev: mongodev
	air --build.cmd "go build -o bins/rest cmd/rest/main.go" --build.bin "./bins/rest"

mongodev: 
	cd docstore && make dev
	
tools: air

air:
	@go install github.com/cosmtrek/air@latest
