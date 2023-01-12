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
ADD ZScaler.crt /usr/local/share/ca-certificates/ZScaler.crt
RUN chmod 644 /usr/local/share/ca-certificates/ZScaler.crt && update-ca-certificates
RUN apt-get update && apt-get install -y ca-certificates git curl netbase wget && rm -rf /var/lib/apt/lists/*

# Copy the source code and create the app directory
COPY . /app
WORKDIR /app

# Compile the Go code
RUN go build -o main .

# Create a new stage for the runtime environment
#FROM ubuntu:latest
FROM ubuntu:latest

# Set the necessary environment variables
ENV MYSQL_ROOT_PASSWORD=Passw0rd,12345
ENV MYSQL_DATABASE=temperature_db
ENV MYSQL_USER=dbuser
ENV MYSQL_PASSWORD=heslo

# CREATE DIRECTORY
RUN mkdir -p /app
# Copy the compiled Go binary and create the app directory
COPY --from=builder /app/main /app/main
# COPY --from=builder /app/create_table.sql /app/create_table.sql
RUN chmod +x /app/main

# Setup working dir
WORKDIR /app

# ENTRYPOINT ./main
CMD ["./main"]
