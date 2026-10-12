package preprocess

import (
	"bytes"
	"image"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/segment"
	"github.com/disintegration/imaging"
)

// Config agrupa los hiperparámetros del pipeline de preprocesamiento.
type Config struct {
	Width      int
	Threshold  int
	BlurRadius float64
}

// Open carga una imagen desde disco.
func Open(path string) (image.Image, error) {
	return imaging.Open(path)
}

// Process redimensiona, convierte a escala de grises, binariza y suaviza
// según la configuración dada.
func Process(src image.Image, cfg Config) image.Image {
	// Procesar con imaging (redimensionar y escala de grises)
	img := imaging.Resize(src, cfg.Width, 0, imaging.Lanczos)
	imgGray := imaging.Grayscale(img)

	// Binarización y suavizado con bild
	imgThresholded := segment.Threshold(imgGray, uint8(cfg.Threshold))
	return blur.Gaussian(imgThresholded, cfg.BlurRadius)
}

// Encode serializa la imagen a JPEG (calidad 95) para alimentar a Tesseract
// sin pasar por disco.
func Encode(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(95))
	return buf.Bytes(), err
}

// Save guarda la imagen en disco, para inspeccionar el preprocesamiento.
func Save(img image.Image, path string) error {
	return imaging.Save(img, path)
}
