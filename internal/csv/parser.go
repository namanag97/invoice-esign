package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// VendorRecord represents one row from the partner CSV.
type VendorRecord struct {
	AccountID          string  `json:"account_id"`
	PartnerName        string  `json:"partner_name"`
	PartnerPhone       string  `json:"partner_phone"`
	PartnerEmail       string  `json:"partner_email"`
	BankAccountNumber  string  `json:"bank_account_number"`
	BankIFSC           string  `json:"bank_ifsc"`
	PartnerAccountType string  `json:"partner_account_type"`
	PartnerCode        string  `json:"partner_code"`
	ARNNo              string  `json:"arn_no"`
	PANNo              string  `json:"pan_no"`
	GSTNo              string  `json:"gst_no"`
	GSTRemarks         string  `json:"gst_remarks"`
	PayoutAmount       float64 `json:"payout_amount"`
	UTR                string  `json:"utr"`
	RowNumber          int     `json:"row_number"`
}

// ParseFile reads a partner CSV and returns vendor records.
func ParseFile(path string) ([]VendorRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading CSV: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("CSV must have a header row and at least one data row")
	}

	headerMap, err := mapHeaders(rows[0])
	if err != nil {
		return nil, err
	}

	var records []VendorRecord
	for i, row := range rows[1:] {
		rowNum := i + 2
		record, err := parseRow(row, headerMap, rowNum)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", rowNum, err)
		}
		records = append(records, record)
	}

	return records, nil
}

var requiredHeaders = []string{
	"partner_name",
	"partner_email",
	"gst_no",
	"payout_amount",
}

func mapHeaders(headers []string) (map[string]int, error) {
	m := make(map[string]int)
	for i, h := range headers {
		normalized := strings.ToLower(strings.TrimSpace(h))
		normalized = strings.ReplaceAll(normalized, " ", "_")
		m[normalized] = i
	}

	for _, required := range requiredHeaders {
		if _, ok := m[required]; !ok {
			return nil, fmt.Errorf("missing required column %q (found: %v)", required, headers)
		}
	}

	return m, nil
}

func parseRow(row []string, headerMap map[string]int, rowNum int) (VendorRecord, error) {
	get := func(col string) string {
		idx, ok := headerMap[col]
		if !ok || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	amountStr := strings.ReplaceAll(get("payout_amount"), ",", "")
	payoutAmount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return VendorRecord{}, fmt.Errorf("invalid payout_amount %q: %w", get("payout_amount"), err)
	}

	return VendorRecord{
		AccountID:          get("account_id"),
		PartnerName:        get("partner_name"),
		PartnerPhone:       get("partner_phone"),
		PartnerEmail:       get("partner_email"),
		BankAccountNumber:  get("bank_account_number"),
		BankIFSC:           get("bank_accountifsccode"),
		PartnerAccountType: get("partner_account_type"),
		PartnerCode:        get("partner_code"),
		ARNNo:              get("arn_no"),
		PANNo:              strings.ToUpper(get("pan_no")),
		GSTNo:              strings.ToUpper(get("gst_no")),
		GSTRemarks:         get("gst_remarks"),
		PayoutAmount:       payoutAmount,
		UTR:                get("utr"),
		RowNumber:          rowNum,
	}, nil
}
