package extract

import (
	"regexp"
	"strconv"
	"strings"
)

// Cambiamos \d{2} por \d{1,2} para que acepte tanto "16.5" como "16.50"
var reTotal = regexp.MustCompile(`(?i)(?:val\.?\s*)?total[^\d]*?(\d+[\.,]\d{1,2})`)

// Total busca el monto total dentro del texto OCR y devuelve su valor
// junto con la línea original donde fue encontrado.
func Total(text string) (amount float64, line string, ok bool) {
	for _, l := range strings.Split(text, "\n") {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		match := reTotal.FindStringSubmatch(l)
		if len(match) <= 1 {
			continue
		}

		v, err := strconv.ParseFloat(strings.ReplaceAll(match[1], ",", "."), 64)
		if err == nil {
			return v, l, true
		}
	}
	return 0, "", false
}
