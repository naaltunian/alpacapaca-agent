image:
	docker build -t paca-agent:latest .

build:
	go build -o build/paca-agent

test:
	go test ./...