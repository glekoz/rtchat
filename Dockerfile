FROM golang:1.20-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /chat ./cmd/chat


FROM alpine:3.18
COPY --from=build /chat /chat
EXPOSE 8080
ENTRYPOINT ["/chat"]
