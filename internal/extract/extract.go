package extract

import (
	"regexp"
	"strconv"
	"strings"
)

// Monto con separador de miles opcional: 1,234.50 | 1.234,50 | 16.50 | 16,5
// El \b final evita capturar "1,23" de "1,234".
const amountPattern = `(\d{1,3}(?:[.,]\d{3})*|\d+)[.,](\d{1,2})\b`

// \btotal\b descarta SUBTOTAL, pero acepta "TOTAL", "Val. Total", "IMPORTE TOTAL".
var reTotal = regexp.MustCompile(`(?i)\btotal\b[^\d]*?` + amountPattern)

// Total busca el monto total dentro del texto OCR. Si hay varias líneas
// candidatas, devuelve la de mayor monto junto con su línea original.
func Total(text string) (amount float64, line string, ok bool) {
	for _, l := range strings.Split(text, "\n") {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		m := reTotal.FindStringSubmatch(l)
		if m == nil {
			continue
		}

		// m[1] = parte entera (puede traer separadores de miles), m[2] = decimales
		whole := strings.NewReplacer(".", "", ",", "").Replace(m[1])
		v, err := strconv.ParseFloat(whole+"."+m[2], 64)
		if err != nil {
			continue
		}

		if !ok || v > amount {
			amount, line, ok = v, l, true
		}
	}
	return amount, line, ok
}
