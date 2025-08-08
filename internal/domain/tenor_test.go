package domain_test

import (
	"testing"

	"github.com/elokanugrah/go-financing-btpns/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestTenor_Calculate(t *testing.T) {
	tests := []struct {
		name             string
		tenorValue       int
		amount           float64
		marginRate       float64
		wantMargin       float64
		wantTotalPayment float64
		wantInstallment  float64
	}{
		{
			name:             "Tenor 6 bulan, margin 20%",
			tenorValue:       6,
			amount:           10000000,
			marginRate:       0.2,
			wantMargin:       1000000,    // (10jt * 0.2 * 6) / 12
			wantTotalPayment: 11000000,   // 10jt + margin
			wantInstallment:  1833333.33, // total / tenor
		},
		{
			name:             "Tenor 12 bulan, margin 15%",
			tenorValue:       12,
			amount:           5000000,
			marginRate:       0.15,
			wantMargin:       750000, // (5jt * 0.15 * 12) / 12
			wantTotalPayment: 5750000,
			wantInstallment:  479166.67,
		},
		{
			name:             "Tenor 24 bulan, margin 10%",
			tenorValue:       24,
			amount:           20000000,
			marginRate:       0.1,
			wantMargin:       4000000, // (20jt * 0.1 * 24) / 12
			wantTotalPayment: 24000000,
			wantInstallment:  1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tenor := domain.Tenor{TenorValue: tt.tenorValue}

			gotInstallment, gotMargin, gotTotal := tenor.Calculate(tt.amount, tt.marginRate)

			assert.InDelta(t, tt.wantMargin, gotMargin, 0.01, "Margin tidak sesuai")
			assert.InDelta(t, tt.wantTotalPayment, gotTotal, 0.01, "Total Payment tidak sesuai")
			assert.InDelta(t, tt.wantInstallment, gotInstallment, 0.01, "Installment tidak sesuai")
		})
	}
}
