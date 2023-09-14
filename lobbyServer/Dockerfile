# syntax=docker/dockerfile:1

FROM golang:1.20.6-alpine

# Set destination for COPY
WORKDIR /appLobbyServer

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./lobbyServer

# Run
CMD ["./lobbyServer"]