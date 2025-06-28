# Etapa de compilación
FROM golang:1.24.1-alpine AS builder
WORKDIR /build

# Instalar dependencias necesarias
RUN apk add --no-cache git

# Copiar todos los archivos necesarios, incluyendo la carpeta pkg
COPY ./server.go .
COPY ./pkg ./pkg
COPY ./playground_files ./playground_files

# Inicializar el módulo Go
RUN go mod init github.com/luis198755/go_playGround_plus/docker

# Agregar dependencias requeridas
RUN go get go.uber.org/zap
RUN go get github.com/pkg/errors
RUN go get github.com/rs/cors

# Instalar todas las dependencias restantes
RUN go mod tidy

# Compilar la aplicación
RUN go build -o server

# Etapa final
FROM golang:1.24.1-alpine
WORKDIR /app

# Copiar archivos necesarios
COPY ./playground_files ./playground_files
COPY --from=builder /build/server .

# Crear directorios necesarios para la aplicación
RUN mkdir -p /app/build

# Crear un usuario no privilegiado
RUN adduser -D appuser && \
    chown -R appuser:appuser /app

# Crear directorios temporales con permisos adecuados
RUN mkdir -p /tmp/go-playground && \
    chmod 777 /tmp/go-playground && \
    mkdir -p /home/appuser/.cache && \
    chmod 777 /home/appuser/.cache

# Cambiar al usuario no privilegiado
USER appuser

# Exponer el puerto configurado por variable de entorno
EXPOSE 8080

# Iniciar el servidor
CMD ["./server"]