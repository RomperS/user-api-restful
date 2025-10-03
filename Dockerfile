# Etapa 1: build
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar binario
RUN go build -o user-api ./cmd

# Etapa 2: runtime
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copiar binario desde el builder
COPY --from=builder /app/user-api .

# Puerto expuesto
EXPOSE 8080

CMD ["./user-api"]
