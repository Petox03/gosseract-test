package ocr

import (
	"fmt"

	"github.com/otiai10/gosseract/v2"
)

// Text ejecuta OCR sobre los bytes de una imagen (PNG/JPEG) con el idioma y el
// modo de segmentación de página dados, y devuelve el texto extraído.
func Text(imageData []byte, lang string, psm gosseract.PageSegMode) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	if err := client.SetImageFromBytes(imageData); err != nil {
		return "", fmt.Errorf("cargar la imagen: %w", err)
	}

	client.SetPageSegMode(psm)
	if err := client.SetLanguage(lang); err != nil {
		return "", fmt.Errorf("configurar idioma %q: %w", lang, err)
	}

	return client.Text()
}
