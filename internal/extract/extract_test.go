package extract

import "testing"

func TestTotal(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		want   float64
		wantOK bool
	}{
		{"simple", "TOTAL $16.50", 16.50, true},
		{"decimal comma", "TOTAL 16,50", 16.50, true},
		{"one decimal", "TOTAL 16.5", 16.5, true},
		{"thousands with comma", "TOTAL $1,234.50", 1234.50, true},
		{"thousands with dot", "TOTAL 1.234,50", 1234.50, true},
		{"ignores subtotal", "SUBTOTAL 100.00\nTOTAL 116.00", 116.00, true},
		{"only subtotal", "SUBTOTAL 100.00", 0, false},
		{"prefers largest amount", "TOTAL IVA 16.00\nTOTAL 116.00", 116.00, true},
		{"no amount", "TOTAL DE ARTICULOS: 3", 0, false},
		{"empty", "", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, ok := Total(tt.text)
			if ok != tt.wantOK || got != tt.want {
				t.Errorf("Total() = (%v, %v), want (%v, %v)", got, ok, tt.want, tt.wantOK)
			}
		})
	}
}
