FROM golang:1.17 as build
WORKDIR /app
COPY go.mod go.sum /app/
COPY src/cmd/go.mod /app/src/cmd/go.mod
RUN go mod download -x
# copy all after mod downloads for Docker Caching Best Practices 
COPY . .
# CGO_ENABLED=0 from https://stackoverflow.com/questions/55106186
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /gobol src/main.go

FROM scratch as scratch
WORKDIR /
COPY --from=build /gobol /gobol
ENV DBHOST="db" \
    DBPORT="3322" \
    TN3270DIR="/app/transactions"
EXPOSE 23567
ENTRYPOINT ["/gobol", "tn3270e"]
