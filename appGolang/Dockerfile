FROM golang:1.19.0-alpine3.16
WORKDIR /appGolang
COPY . .
RUN go mod download
RUN go mod tidy
EXPOSE 8000
CMD ["go", "run", "main.go"]