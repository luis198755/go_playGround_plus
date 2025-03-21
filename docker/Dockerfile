# Etapa de compilación
FROM golang:1.24.1-alpine AS builder
WORKDIR /build
COPY ./server.go .
RUN go mod init server
RUN go mod tidy
RUN go build -o server

# Etapa final
FROM golang:1.24.1-alpine
WORKDIR /app
COPY ./playground_files ./playground_files
COPY --from=builder /build/server .

# Crear un usuario no privilegiado
RUN adduser -D appuser
USER appuser

# Hacer el sistema de archivos de solo lectura
VOLUME ["/app"]
#EXPOSE 8080
CMD ["./server"]