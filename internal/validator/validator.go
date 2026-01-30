package validator

import (
	"fmt"
	"regexp"
	"strings"

	csvpkg "github.com/namanag97/invoice-esign/internal/csv"
	"github.com/namanag97/invoice-esign/internal/gst"
)

var panPattern = regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]$`)

var phonePattern = regexp.MustCompile(`^\d{10}$`)

type ValidationResult struct {
	Record       csvpkg.VendorRecord
	IsValid      bool
	Issues       []string
	CleanedPhone string
	CleanedGSTIN string
	CleanedPAN   string
	VendorState  string
}

// ValidateRecord checks a vendor record and returns validation results.
func ValidateRecord(r csvpkg.VendorRecord) ValidationResult {
	result := ValidationResult{
		Record:  r,
		IsValid: true,
	}

	// Name
	if strings.TrimSpace(r.PartnerName) == "" {
		result.Issues = append(result.Issues, "partner_name is empty")
		result.IsValid = false
	}

	// Email
	if strings.TrimSpace(r.PartnerEmail) == "" {
		result.Issues = append(result.Issues, "partner_email is empty")
		result.IsValid = false
	} else if !strings.Contains(r.PartnerEmail, "@") {
		result.Issues = append(result.Issues, fmt.Sprintf("partner_email %q is not a valid email", r.PartnerEmail))
		result.IsValid = false
	}

	// Payout amount
	if r.PayoutAmount <= 0 {
		result.Issues = append(result.Issues, fmt.Sprintf("payout_amount must be positive, got %.2f", r.PayoutAmount))
		result.IsValid = false
	}

	// GSTIN
	cleanedGSTIN := strings.TrimSpace(strings.ToUpper(r.GSTNo))
	if cleanedGSTIN == "" {
		result.Issues = append(result.Issues, "gst_no is empty")
		result.IsValid = false
	} else if err := gst.ValidateGSTIN(cleanedGSTIN); err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("gst_no: %v", err))
		result.IsValid = false
	} else {
		result.CleanedGSTIN = cleanedGSTIN
		result.VendorState = gst.StateFromGSTIN(cleanedGSTIN)
	}

	// PAN
	cleanedPAN := strings.TrimSpace(strings.ToUpper(r.PANNo))
	if cleanedPAN != "" {
		if !panPattern.MatchString(cleanedPAN) {
			result.Issues = append(result.Issues, fmt.Sprintf("pan_no %q is not valid format", cleanedPAN))
		} else {
			result.CleanedPAN = cleanedPAN
		}
	}

	// Phone
	result.CleanedPhone = cleanPhone(r.PartnerPhone)

	return result
}

func cleanPhone(phone string) string {
	// Strip everything except digits
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	// Remove country code
	if strings.HasPrefix(digits, "91") && len(digits) == 12 {
		digits = digits[2:]
	} else if strings.HasPrefix(digits, "0") && len(digits) == 11 {
		digits = digits[1:]
	}

	if phonePattern.MatchString(digits) {
		return digits
	}
	return phone // return original if can't clean
}
