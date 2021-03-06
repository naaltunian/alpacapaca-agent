image:
	docker build -t paca-agent:latest .

test:
	go test ./...