#!/bin/bash

echo "🚀 Configurando el nodo de la blockchain DAG..."

# 1. Instalar Go si no está instalado
if ! command -v go &> /dev/null
then
    echo "⚠️  Go no está instalado. Instalándolo..."
    wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo "✅ Go instalado correctamente"
fi

# 2. Configurar variables de entorno
export NODE_PORT=8080
export P2P_PORT=6000
export DATA_DIR="$HOME/.fernet-token"
mkdir -p $DATA_DIR

echo "📂 Carpeta de datos creada en $DATA_DIR"

# 3. Descargar dependencias del proyecto
echo "📦 Descargando dependencias..."
go mod tidy

# 4. Generar clave privada para el nodo
if [ ! -f "$DATA_DIR/private_key.pem" ]; then
    echo "🔑 Generando clave privada del nodo..."
    openssl genpkey -algorithm RSA -out "$DATA_DIR/private_key.pem"
fi

# 5. Iniciar el nodo
echo "🚀 Iniciando nodo en el puerto $NODE_PORT..."
go run cmd/server/main.go

echo "✅ Nodo ejecutándose correctamente. ¡Bienvenido a la red DAG de Fernet!"
