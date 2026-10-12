# gosseract-test

[![CI](https://github.com/Petox03/gosseract-test/actions/workflows/ci.yml/badge.svg)](https://github.com/Petox03/gosseract-test/actions/workflows/ci.yml)

Lab de OCR en Go para extraer el **total** de tickets con Tesseract:
preprocesamiento de imagen + OCR + extracción por regex, hecho a mano
y self-hosted en Docker.

## Uso

```bash
docker compose up -d
docker compose exec app go run ./cmd/ocr <ticket>.png
```

Las imágenes se leen desde `receipts/` (ignorada en git — agrega las tuyas).

Flags para el pipeline (todas opcionales, `-h` muestra los defaults):

- `-out` — guarda la imagen preprocesada en disco para inspeccionarla
- `-dir` — carpeta base de las imágenes (default `receipts`)
- `-lang` — idioma de Tesseract (`spa`, `eng`, …)
- `-psm` — modo de segmentación de página (cómo Tesseract interpreta el layout)
- `-width` — ancho en px al redimensionar (`0` conserva el original)
- `-threshold` — umbral de binarización (0-255)
- `-blur` — radio del desenfoque gaussiano (`0` lo desactiva)

## Estructura

- `cmd/ocr` — wiring: open → preprocess → OCR → extract
- `internal/preprocess` — resize, grayscale, binarize, blur, encode
- `internal/ocr` — wrapper de gosseract (imagen en memoria, sin temporales)
- `internal/extract` — extracción del total por regex
