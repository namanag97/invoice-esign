package template

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/namanag97/invoice-esign/internal/config"
	"github.com/namanag97/invoice-esign/internal/gst"
	"github.com/namanag97/invoice-esign/internal/validator"
)

// InvoiceData holds all data needed to render an invoice.
type InvoiceData struct {
	// Vendor
	VendorName  string
	VendorState string
	VendorGSTIN string
	VendorPhone string
	VendorEmail string

	// Invoice metadata
	InvoiceNo   string
	InvoiceDate string
	Terms       string
	DueDate     string

	// Company (Bill To)
	CompanyName    string
	CompanyGSTIN   string
	CompanyAddress string
	PlaceOfSupply  string

	// Amounts
	BaseAmount   string
	CGSTAmt      string
	SGSTAmt      string
	IGSTAmt      string
	TotalWithTax string
	Adjustment   string
	BalanceDue   string
	AmountWords  string
	TaxNote      string

	// Bank details
	BankName       string
	BankAccount    string
	PaymentRef     string

	// Display control for CGST/SGST vs IGST columns
	CGSTSGSTStyle template.HTMLAttr
	IGSTStyle     template.HTMLAttr
}

// RenderInvoice generates an HTML invoice for a validated vendor record.
func RenderInvoice(cfg *config.Config, v validator.ValidationResult, invoiceNo, invoiceDate, dueDate string) (string, error) {
	interState := !gst.IsSameState(v.CleanedGSTIN, cfg.CompanyGSTIN)
	breakdown := gst.CalculateGST(v.Record.PayoutAmount, 18.0, interState)

	cgstSGSTStyle := template.HTMLAttr("")
	igstStyle := template.HTMLAttr(`style="display: none;"`)
	if interState {
		cgstSGSTStyle = template.HTMLAttr(`style="display: none;"`)
		igstStyle = template.HTMLAttr("")
	}

	taxNote := "CGST @9% + SGST @9%"
	if interState {
		taxNote = "IGST @18%"
	}

	state := v.VendorState
	if state == "" {
		state = gst.StateFromGSTIN(v.CleanedGSTIN)
	}

	data := InvoiceData{
		VendorName:     v.Record.PartnerName,
		VendorState:    state,
		VendorGSTIN:    v.CleanedGSTIN,
		VendorPhone:    v.CleanedPhone,
		VendorEmail:    v.Record.PartnerEmail,
		InvoiceNo:      invoiceNo,
		InvoiceDate:    invoiceDate,
		Terms:          fmt.Sprintf("Net %d", cfg.PaymentTermsDays),
		DueDate:        dueDate,
		CompanyName:    cfg.CompanyName,
		CompanyGSTIN:   cfg.CompanyGSTIN,
		CompanyAddress: cfg.CompanyAddress,
		PlaceOfSupply:  fmt.Sprintf("%s (%s)", cfg.CompanyState, cfg.CompanyStateCode()),
		BaseAmount:     formatINR(breakdown.TaxableAmount),
		CGSTAmt:        formatINR(breakdown.CGSTAmount),
		SGSTAmt:        formatINR(breakdown.SGSTAmount),
		IGSTAmt:        formatINR(breakdown.IGSTAmount),
		TotalWithTax:   formatINR(breakdown.TotalAmount),
		Adjustment:     "₹0.00",
		BalanceDue:     formatINR(breakdown.TotalAmount),
		AmountWords:    amountToWordsINR(breakdown.TotalAmount),
		TaxNote:        taxNote,
		BankName:       bankNameFromIFSC(v.Record.BankIFSC),
		BankAccount:    maskAccount(v.Record.BankAccountNumber),
		PaymentRef:     v.Record.UTR,
		CGSTSGSTStyle:  cgstSGSTStyle,
		IGSTStyle:      igstStyle,
	}

	if data.PaymentRef == "" {
		data.PaymentRef = "Pending"
	}

	tmplPath := cfg.TemplatePath
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return "", fmt.Errorf("reading template %s: %w", tmplPath, err)
	}

	tmpl, err := template.New("invoice").Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

// SaveHTML writes the rendered HTML to a file and returns the path.
func SaveHTML(html, outputDir, invoiceNo string) (string, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", err
	}
	safe := strings.ReplaceAll(invoiceNo, "/", "_")
	safe = strings.ReplaceAll(safe, "\\", "_")
	outPath := filepath.Join(outputDir, safe+".html")
	return outPath, os.WriteFile(outPath, []byte(html), 0o644)
}

func formatINR(amount float64) string {
	// Indian number formatting: ₹1,23,456.78
	negative := amount < 0
	amount = math.Abs(amount)

	intPart := int64(amount)
	decPart := int64(math.Round((amount - float64(intPart)) * 100))

	s := fmt.Sprintf("%d", intPart)
	// Indian grouping: last 3 digits, then groups of 2
	if len(s) > 3 {
		result := s[len(s)-3:]
		s = s[:len(s)-3]
		for len(s) > 2 {
			result = s[len(s)-2:] + "," + result
			s = s[:len(s)-2]
		}
		if len(s) > 0 {
			result = s + "," + result
		}
		s = result
	}

	prefix := "₹"
	if negative {
		prefix = "-₹"
	}
	return fmt.Sprintf("%s%s.%02d", prefix, s, decPart)
}

func maskAccount(account string) string {
	if len(account) <= 4 {
		return "XXXX"
	}
	return "XXXX-" + account[len(account)-4:]
}

func bankNameFromIFSC(ifsc string) string {
	if len(ifsc) < 4 {
		return "Unknown"
	}
	prefix := strings.ToUpper(ifsc[:4])
	banks := map[string]string{
		"HDFC": "HDFC Bank",
		"ICIC": "ICICI Bank",
		"SBIN": "State Bank of India",
		"BKID": "Bank of India",
		"YESB": "YES Bank",
		"KKBK": "Kotak Mahindra Bank",
		"UTIB": "Axis Bank",
		"PUNB": "Punjab National Bank",
		"CNRB": "Canara Bank",
		"UBIN": "Union Bank of India",
		"IDFB": "IDFC First Bank",
		"BARB": "Bank of Baroda",
	}
	if name, ok := banks[prefix]; ok {
		return name
	}
	return prefix + " Bank"
}

// amountToWordsINR converts a float to Indian rupees in words.
// Simplified implementation for common amounts.
func amountToWordsINR(amount float64) string {
	rupees := int64(amount)
	paise := int64(math.Round((amount - float64(rupees)) * 100))

	result := numberToWords(rupees) + " Rupees"
	if paise > 0 {
		result += " and " + numberToWords(paise) + " Paise"
	}
	result += " Only"
	return result
}

func numberToWords(n int64) string {
	if n == 0 {
		return "Zero"
	}

	ones := []string{"", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine",
		"Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"}
	tens := []string{"", "", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"}

	// Indian system: Lakh (1,00,000) and Crore (1,00,00,000)
	var parts []string

	if n >= 10000000 {
		parts = append(parts, numberToWords(n/10000000)+" Crore")
		n %= 10000000
	}
	if n >= 100000 {
		parts = append(parts, numberToWords(n/100000)+" Lakh")
		n %= 100000
	}
	if n >= 1000 {
		parts = append(parts, numberToWords(n/1000)+" Thousand")
		n %= 1000
	}
	if n >= 100 {
		parts = append(parts, numberToWords(n/100)+" Hundred")
		n %= 100
	}
	if n >= 20 {
		part := tens[n/10]
		if n%10 != 0 {
			part += " " + ones[n%10]
		}
		parts = append(parts, part)
	} else if n > 0 {
		parts = append(parts, ones[n])
	}

	return strings.Join(parts, " ")
}
