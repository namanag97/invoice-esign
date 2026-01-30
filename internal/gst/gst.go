package gst

import (
	"fmt"
	"regexp"
	"strings"
)

// Indian GST number format: 22AAAAA0000A1Z5
// - Digits 1-2: State code (01-38)
// - Digits 3-12: PAN number
// - Digit 13: Entity number
// - Digit 14: 'Z' (default)
// - Digit 15: Check digit

var gstinPattern = regexp.MustCompile(`^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z][0-9A-Z]Z[0-9A-Z]$`)

// State codes as per Indian GST.
var stateCodes = map[string]string{
	"01": "Jammu & Kashmir",
	"02": "Himachal Pradesh",
	"03": "Punjab",
	"04": "Chandigarh",
	"05": "Uttarakhand",
	"06": "Haryana",
	"07": "Delhi",
	"08": "Rajasthan",
	"09": "Uttar Pradesh",
	"10": "Bihar",
	"11": "Sikkim",
	"12": "Arunachal Pradesh",
	"13": "Nagaland",
	"14": "Manipur",
	"15": "Mizoram",
	"16": "Tripura",
	"17": "Meghalaya",
	"18": "Assam",
	"19": "West Bengal",
	"20": "Jharkhand",
	"21": "Odisha",
	"22": "Chhattisgarh",
	"23": "Madhya Pradesh",
	"24": "Gujarat",
	"25": "Daman & Diu",
	"26": "Dadra & Nagar Haveli",
	"27": "Maharashtra",
	"28": "Andhra Pradesh",
	"29": "Karnataka",
	"30": "Goa",
	"31": "Lakshadweep",
	"32": "Kerala",
	"33": "Tamil Nadu",
	"34": "Puducherry",
	"35": "Andaman & Nicobar Islands",
	"36": "Telangana",
	"37": "Andhra Pradesh (New)",
	"38": "Ladakh",
}

// ValidateGSTIN checks if a GSTIN is valid and returns the state name.
func ValidateGSTIN(gstin string) error {
	gstin = strings.TrimSpace(strings.ToUpper(gstin))
	if len(gstin) != 15 {
		return fmt.Errorf("GSTIN must be 15 characters, got %d", len(gstin))
	}
	if !gstinPattern.MatchString(gstin) {
		return fmt.Errorf("GSTIN %q does not match expected format", gstin)
	}
	code := gstin[:2]
	if _, ok := stateCodes[code]; !ok {
		return fmt.Errorf("invalid state code %q in GSTIN", code)
	}
	return nil
}

// StateFromGSTIN extracts the state name from a GSTIN.
func StateFromGSTIN(gstin string) string {
	gstin = strings.TrimSpace(strings.ToUpper(gstin))
	if len(gstin) < 2 {
		return "Unknown"
	}
	state, ok := stateCodes[gstin[:2]]
	if !ok {
		return "Unknown"
	}
	return state
}

// PANFromGSTIN extracts the PAN number from a GSTIN.
func PANFromGSTIN(gstin string) string {
	gstin = strings.TrimSpace(strings.ToUpper(gstin))
	if len(gstin) < 12 {
		return ""
	}
	return gstin[2:12]
}

// IsSameState checks if two GSTINs belong to the same state.
// Same state → CGST + SGST. Different state → IGST.
func IsSameState(gstin1, gstin2 string) bool {
	if len(gstin1) < 2 || len(gstin2) < 2 {
		return false
	}
	return strings.ToUpper(gstin1[:2]) == strings.ToUpper(gstin2[:2])
}

// GSTBreakdown holds the tax calculation for an invoice line.
type GSTBreakdown struct {
	TaxableAmount float64 `json:"taxable_amount"`
	GSTRate       float64 `json:"gst_rate"`
	CGSTRate      float64 `json:"cgst_rate"`
	CGSTAmount    float64 `json:"cgst_amount"`
	SGSTRate      float64 `json:"sgst_rate"`
	SGSTAmount    float64 `json:"sgst_amount"`
	IGSTRate      float64 `json:"igst_rate"`
	IGSTAmount    float64 `json:"igst_amount"`
	TotalTax      float64 `json:"total_tax"`
	TotalAmount   float64 `json:"total_amount"`
	IsInterState  bool    `json:"is_inter_state"`
}

// CalculateGST computes the GST breakdown.
// If inter-state (different states), IGST applies. Otherwise CGST + SGST (split equally).
func CalculateGST(taxableAmount, gstRatePercent float64, interState bool) GSTBreakdown {
	b := GSTBreakdown{
		TaxableAmount: taxableAmount,
		GSTRate:       gstRatePercent,
		IsInterState:  interState,
	}

	totalTax := taxableAmount * gstRatePercent / 100

	if interState {
		b.IGSTRate = gstRatePercent
		b.IGSTAmount = round2(totalTax)
	} else {
		b.CGSTRate = gstRatePercent / 2
		b.CGSTAmount = round2(totalTax / 2)
		b.SGSTRate = gstRatePercent / 2
		b.SGSTAmount = round2(totalTax / 2)
	}

	b.TotalTax = round2(totalTax)
	b.TotalAmount = round2(taxableAmount + totalTax)
	return b
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
