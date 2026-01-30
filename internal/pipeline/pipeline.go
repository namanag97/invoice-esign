package pipeline

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/namanag97/invoice-esign/internal/config"
	csvpkg "github.com/namanag97/invoice-esign/internal/csv"
	"github.com/namanag97/invoice-esign/internal/gst"
	"github.com/namanag97/invoice-esign/internal/mailer"
	tmplpkg "github.com/namanag97/invoice-esign/internal/template"
	"github.com/namanag97/invoice-esign/internal/validator"
)

type Pipeline struct {
	cfg      *config.Config
	csvPath  string
}

func New(cfg *config.Config, csvPath string) *Pipeline {
	return &Pipeline{cfg: cfg, csvPath: csvPath}
}

// Run executes the full pipeline: parse → validate → generate → email.
func (p *Pipeline) Run() error {
	fmt.Println("=== invoice-esign pipeline ===")

	// Stage 1: Parse CSV
	fmt.Printf("\n[1/4] Parsing CSV: %s\n", p.csvPath)
	records, err := csvpkg.ParseFile(p.csvPath)
	if err != nil {
		return fmt.Errorf("CSV parse: %w", err)
	}
	fmt.Printf("  Found %d vendor record(s)\n", len(records))

	// Stage 2: Validate
	fmt.Println("\n[2/4] Validating records...")
	var valid []validator.ValidationResult
	var invalid []validator.ValidationResult
	for _, r := range records {
		result := validator.ValidateRecord(r)
		if result.IsValid {
			valid = append(valid, result)
		} else {
			invalid = append(invalid, result)
			fmt.Printf("  ✗ Row %d (%s): %s\n", r.RowNumber, r.PartnerName, strings.Join(result.Issues, "; "))
		}
	}
	fmt.Printf("  Valid: %d, Invalid: %d\n", len(valid), len(invalid))

	if len(valid) == 0 {
		return fmt.Errorf("no valid records to process")
	}

	// Stage 3: Generate invoices
	fmt.Println("\n[3/4] Generating invoices...")
	now := time.Now()
	invoiceDate := now.Format("02-01-2006")
	dueDate := now.AddDate(0, 0, p.cfg.PaymentTermsDays).Format("02-01-2006")
	invoiceDir := filepath.Join(p.cfg.OutputDir, "invoices")

	type generated struct {
		result    validator.ValidationResult
		htmlPath  string
		invoiceNo string
		total     string
	}
	var invoices []generated

	for i, v := range valid {
		// Generate invoice number
		partnerShort := strings.ToUpper(v.Record.PartnerCode)
		if partnerShort == "" {
			partnerShort = fmt.Sprintf("V%03d", i+1)
		}
		monthCode := now.Format("0106")
		invoiceNo := fmt.Sprintf("VM-%s-%s-%03d", partnerShort, monthCode, i+1)

		html, err := tmplpkg.RenderInvoice(p.cfg, v, invoiceNo, invoiceDate, dueDate)
		if err != nil {
			fmt.Printf("  ✗ %s: render failed: %v\n", v.Record.PartnerName, err)
			continue
		}

		htmlPath, err := tmplpkg.SaveHTML(html, invoiceDir, invoiceNo)
		if err != nil {
			fmt.Printf("  ✗ %s: save failed: %v\n", v.Record.PartnerName, err)
			continue
		}

		interState := !gst.IsSameState(v.CleanedGSTIN, p.cfg.CompanyGSTIN)
		breakdown := gst.CalculateGST(v.Record.PayoutAmount, 18.0, interState)

		invoices = append(invoices, generated{
			result:    v,
			htmlPath:  htmlPath,
			invoiceNo: invoiceNo,
			total:     fmt.Sprintf("₹%.2f", breakdown.TotalAmount),
		})
		fmt.Printf("  ✓ %s → %s\n", v.Record.PartnerName, htmlPath)
	}

	fmt.Printf("  Generated %d invoice(s)\n", len(invoices))

	// Stage 4: Send emails
	fmt.Println("\n[4/4] Dispatching emails...")
	mailClient := mailer.New(p.cfg)
	sent := 0
	failed := 0

	for _, inv := range invoices {
		err := mailClient.Send(mailer.EmailPayload{
			ToEmail:        inv.result.Record.PartnerEmail,
			ToName:         inv.result.Record.PartnerName,
			Subject:        fmt.Sprintf("Tax Invoice %s - %s", inv.invoiceNo, p.cfg.CompanyName),
			InvoiceNumber:  inv.invoiceNo,
			InvoiceDate:    invoiceDate,
			DueDate:        dueDate,
			TotalAmount:    inv.total,
			AttachmentPath: inv.htmlPath,
		})
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", inv.result.Record.PartnerEmail, err)
			failed++
		} else {
			fmt.Printf("  ✓ %s → %s\n", inv.invoiceNo, inv.result.Record.PartnerEmail)
			sent++
		}
	}

	fmt.Printf("\n=== Done. Invoices: %d generated, Emails: %d sent, %d failed ===\n", len(invoices), sent, failed)
	return nil
}

// ValidateOnly runs just the validation stage and prints a report.
func (p *Pipeline) ValidateOnly() error {
	records, err := csvpkg.ParseFile(p.csvPath)
	if err != nil {
		return err
	}

	fmt.Printf("Validating %d record(s)...\n\n", len(records))
	validCount := 0
	for _, r := range records {
		result := validator.ValidateRecord(r)
		if result.IsValid {
			fmt.Printf("  ✓ Row %d: %s (%s) — %s, payout ₹%.2f\n",
				r.RowNumber, r.PartnerName, result.CleanedGSTIN, result.VendorState, r.PayoutAmount)
			validCount++
		} else {
			fmt.Printf("  ✗ Row %d: %s — %s\n",
				r.RowNumber, r.PartnerName, strings.Join(result.Issues, "; "))
		}
	}
	fmt.Printf("\nTotal: %d, Valid: %d, Invalid: %d\n", len(records), validCount, len(records)-validCount)
	return nil
}

// GenerateOnly runs parse → validate → generate (no email).
func (p *Pipeline) GenerateOnly() error {
	records, err := csvpkg.ParseFile(p.csvPath)
	if err != nil {
		return err
	}

	now := time.Now()
	invoiceDate := now.Format("02-01-2006")
	dueDate := now.AddDate(0, 0, p.cfg.PaymentTermsDays).Format("02-01-2006")
	invoiceDir := filepath.Join(p.cfg.OutputDir, "invoices")

	generated := 0
	for i, r := range records {
		v := validator.ValidateRecord(r)
		if !v.IsValid {
			fmt.Printf("  ✗ Row %d: %s — skipped (%s)\n", r.RowNumber, r.PartnerName, strings.Join(v.Issues, "; "))
			continue
		}

		partnerShort := strings.ToUpper(r.PartnerCode)
		if partnerShort == "" {
			partnerShort = fmt.Sprintf("V%03d", i+1)
		}
		monthCode := now.Format("0106")
		invoiceNo := fmt.Sprintf("VM-%s-%s-%03d", partnerShort, monthCode, i+1)

		html, err := tmplpkg.RenderInvoice(p.cfg, v, invoiceNo, invoiceDate, dueDate)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", r.PartnerName, err)
			continue
		}

		path, err := tmplpkg.SaveHTML(html, invoiceDir, invoiceNo)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", r.PartnerName, err)
			continue
		}

		fmt.Printf("  ✓ %s → %s\n", r.PartnerName, path)
		generated++
	}

	fmt.Printf("\nGenerated %d invoice(s) in %s\n", generated, invoiceDir)
	return nil
}
