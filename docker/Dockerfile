# Etapa de compilación
FROM golang:1.24.1-alpine AS builder
WORKDIR /build
# Copiar todos los archivos necesarios, incluyendo la carpeta pkg
COPY ./server.go .
COPY ./pkg ./pkg
# Inicializar el módulo y sus dependencias
RUN go mod init github.com/luis198755/go_playGround_plus/docker
RUN go mod tidy
# Compilar la aplicación
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