# Erstellen Sie ein temporäres Build-Image
FROM golang:1.21.1 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./lobbyServer

# Erstellen Sie das endgültige Image auf Basis von Distroless
FROM gcr.io/distroless/base

COPY --from=build /app/lobbyServer /app/lobbyServer

WORKDIR /app

CMD ["./lobbyServer"]
