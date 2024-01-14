##############################################################
# Multistage 
# Compile Go API First and put it to UBUNTU Container together
# with MySQL Database
##############################################################
FROM golang:latest as builder

# Set the necessary environment variables
ENV MYSQL_ROOT_PASSWORD=Passw0rd,12345
ENV MYSQL_DATABASE=temperature_db
ENV MYSQL_USER=dbuser
ENV MYSQL_PASSWORD=heslo

# Install Certificate Required in OFFICE (BECAUSE OF FUCKING MAN IN THE MIDDLE calle ZScaler)
# ADD ZScaler.crt /usr/local/share/ca-certificates/ZScaler.crt
# RUN chmod 644 /usr/local/share/ca-certificates/ZScaler.crt && update-ca-certificates
RUN apt-get update && apt-get install -y git curl

# Copy the source code and create the app directory
COPY . /app
WORKDIR /app

# Compile the Go code
RUN go build -o main .

# Create a new stage for the runtime environment
FROM ubuntu:latest

# Install MySQL
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y mysql-server

# Set the necessary environment variables
ENV MYSQL_ROOT_PASSWORD=Passw0rd,12345
ENV MYSQL_DATABASE=temperature_db
ENV MYSQL_USER=dbuser
ENV MYSQL_PASSWORD=heslo

# Copy the compiled Go binary and create the app directory
COPY --from=builder /app/main /app/main
COPY --from=builder /app/create_table.sql /app/create_table.sql
WORKDIR /app

# Setup home directory
RUN usermod -d /var/lib/mysql/ mysql

# Create the database and table
RUN service mysql start && \
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -e "CREATE DATABASE $MYSQL_DATABASE" && \
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -e "CREATE USER '$MYSQL_USER'@'%' IDENTIFIED BY '$MYSQL_PASSWORD'" && \
    mysql -uroot -p$MYSQL_ROOT_PASSWORD -e "GRANT ALL ON $MYSQL_DATABASE.* TO '$MYSQL_USER'@'%'" && \
    mysql -uroot -p$MYSQL_ROOT_PASSWORD $MYSQL_DATABASE < create_table.sql

# Run the MySQL server and the compiled Go binary when the container starts
CMD service mysql start && ./main
