FROM ubuntu:20.04

# Install Go and MySQL
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y git golang-go mysql-server


# Set the necessary environment variables
ENV MYSQL_ROOT_PASSWORD=Passw0rd,12345
ENV MYSQL_DATABASE=temperature_db
ENV MYSQL_USER=dbuser
ENV MYSQL_PASSWORD=heslo 


# Create the app directory and set it as the working directory
RUN mkdir -p /app


# Copy the Go code and compiled binary into the container
COPY . /app

# Setup Workdir
WORKDIR /app

# Setup home directory
RUN usermod -d /var/lib/mysql/ mysql

# Compile the Go code
RUN go build -o main .

# Create the database and table
RUN service mysql start && \
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -e "CREATE DATABASE $MYSQL_DATABASE" && \
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -e "CREATE USER '$MYSQL_USER'@'%' IDENTIFIED BY '$MYSQL_PASSWORD'" && \
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -e "GRANT ALL ON $MYSQL_DATABASE.* TO '$MYSQL_USER'@'%'" && \
    mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE < create_table.sql

# Run the MySQL server and the compiled Go binary when the container starts
CMD service mysql start && ./main 2> log.txt