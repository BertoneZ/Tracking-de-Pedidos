# ETAPA 1: Compilación (Builder)
FROM golang:1.24-alpine AS builder

# Set el directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias y descargamos
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilamos el binario 
RUN go build -o main ./cmd/api/main.go

# ETAPA 2: Ejecución (Runner)
FROM alpine:latest

WORKDIR /root/

# Traemos solo el binario compilado de la etapa anterior
COPY --from=builder /app/main .
# También copiamos el .env si lo necesitás (aunque es mejor por Compose)
#COPY .env . 

# Exponemos el puerto
EXPOSE 8080

# Comando para arrancar
CMD ["./main"]