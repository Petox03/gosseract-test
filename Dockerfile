FROM golang:1.27.2-alpine

WORKDIR /app

# Dependencias del sistema para compilar con CGO y para Tesseract/Leptonica
RUN apk add --no-cache \
    git \
    build-base \
    pkgconf \
    tesseract-ocr \
    tesseract-ocr-dev \
    leptonica-dev \
    tesseract-ocr-data-eng \
    tesseract-ocr-data-spa

# gosseract usa CGO, así que lo dejamos habilitado de forma explícita
ENV CGO_ENABLED=1

# Redirigir cachés de Go a rutas que serán persistidas vía volúmenes nombrados
ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/root/.cache/go-mod

# Carpeta donde Tesseract busca los archivos de idioma (.traineddata)
ENV TESSDATA_PREFIX=/usr/share/tessdata

# Falla en el build (no en ejecución) si faltan los datos del idioma inglés
RUN test -f "$TESSDATA_PREFIX/eng.traineddata" \
    || (echo "ERROR: falta eng.traineddata en $TESSDATA_PREFIX" && exit 1)

# gopls: motor de autocompletado para desarrollo
RUN go install golang.org/x/tools/gopls@v0.23.0

# Copiar go.mod y go.sum (el asterisco evita error si go.sum aún no existe)
COPY go.mod go.sum* ./

# Descargar dependencias a la caché de módulos
RUN go mod download

CMD ["sleep", "infinity"]