FROM golang:1.17 as build
WORKDIR /app
# COPY go.mod .
# COPY go.sum .
COPY . .
RUN go mod download -x
RUN go build -o /gobol src/main.go
ENV DBHOST="db" \
    DBPORT="3322" \
    TN3270DIR="/app/transactions"
EXPOSE 23567
ENTRYPOINT ["/gobol", "tn3270e"]

# FROM debian:bullseye-slim as bullseye-slim
# COPY --from=build /gobol /usr/sbin/gobol
# ENTRYPOINT ["/usr/sbin/gobol"