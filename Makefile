build:
	docker build -t go-api-mysql .

run:
	docker run -d -p 80:8080 go-api-mysql
