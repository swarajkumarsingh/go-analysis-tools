run:
	go run main.go

build:
	go build

install:
	go mod tidy

dev:
	nodemon --exec go run main.go

run_prometheus:
	docker run -d -p 9090:9090 -v ./prometheus.yml:/etc/prometheus prom/prometheus

run_grafana:
	docker run -d --name=grafana -p 3000:3000 grafana/grafana-enterprise
	
config:
	make run_prometheus && make run_grafana