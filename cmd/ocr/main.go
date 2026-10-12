package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Petox03/gosseract-test/internal/extract"
	"github.com/Petox03/gosseract-test/internal/ocr"
	"github.com/Petox03/gosseract-test/internal/preprocess"
	"github.com/otiai10/gosseract/v2"
)

func main() {
	var cfg preprocess.Config

	out := flag.String("out", "", "save the preprocessed image to this path")
	dir := flag.String("dir", "receipts", "base directory for ticket images")
	lang := flag.String("lang", "spa", "tesseract language")
	psm := flag.Int("psm", int(gosseract.PSM_SINGLE_BLOCK), "page segmentation mode (6=single block, 3=auto, 11=sparse)")
	flag.IntVar(&cfg.Width, "width", 1200, "resize width in px (0 keeps original size)")
	flag.IntVar(&cfg.Threshold, "threshold", 128, "binarization threshold (0-255)")
	flag.Float64Var(&cfg.BlurRadius, "blur", 1.0, "gaussian blur radius (0 disables)")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/ocr [flags] <ticket.png>")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	imagePath := filepath.Join(*dir, flag.Arg(0))

	// 1. Abrir la imagen original usando el parámetro recibido
	src, err := preprocess.Open(imagePath)
	if err != nil {
		log.Fatalf("failed to open image %s: %v", imagePath, err)
	}

	// 2. Preprocesar (resize, escala de grises, binarización y suavizado)
	processed := preprocess.Process(src, cfg)

	// 3. Guardar copia del preprocesamiento si se pidió, y codificar a JPEG en memoria
	if *out != "" {
		if err := preprocess.Save(processed, *out); err != nil {
			log.Fatalf("failed to save preprocessed image: %v", err)
		}
	}
	imageData, err := preprocess.Encode(processed)
	if err != nil {
		log.Fatalf("failed to encode image: %v", err)
	}

	// 4. OCR con Gosseract
	text, err := ocr.Text(imageData, *lang, gosseract.PageSegMode(*psm))
	if err != nil {
		log.Fatalf("OCR failed: %v", err)
	}

	fmt.Println("--- TEXTO BRUTO EXTRAÍDO POR OCR ---")
	fmt.Println(text)
	fmt.Println("------------------------------------")

	// 5. Post-procesamiento y extracción de datos
	fmt.Println("\n--- RESULTADOS DEL POST-PROCESAMIENTO ---")
	if valor, linea, ok := extract.Total(text); ok {
		fmt.Printf(" [✅] Total Detectado: %.2f (Línea original: '%s')\n", valor, linea)
	} else {
		fmt.Println(" [⚠️] No se pudo detectar un 'Total' claro mediante expresiones regulares en este ticket.")
	}
}
