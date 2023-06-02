FROM golang:1.18 as build
WORKDIR /app/
COPY go.mod go.sum /app/
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags immudb -ldflags="-w -s" -o /transactions/tran1 tran1.go

FROM scratch 
WORKDIR /
COPY --from=build /transactions/tran1 /transactions/tran1
# COPY  SCREEN.hlasm /transactions/screens/SCREEN.hlasm
ENV DBHOST="db" \
    DBPORT="3322"
ENTRYPOINT ["/transactions/tran1"]