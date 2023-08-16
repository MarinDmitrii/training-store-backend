BINARY=training-store-backend

.PHONY: openapi_http
openapi_http:
	@./scripts/openapi-http.sh order internal/order/ports ports

.PHONY: build
build:
	go build -o ${BINARY} cmd/main.go

.PHONY: clean
clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi

.PHONY: docker
docker:
	docker build -t training-store-backend .

.PHONY: run
run:
	docker-compose up --build -d

.PHONY: stop
stop:
	docker-compose down
