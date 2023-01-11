build:
	docker build -t go-api .

run:
	docker run -v mysql_vol:/var/lib/mysql --name gotempapi-mysql -d -p 80:8080 go-api
