package mailer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/namanag97/invoice-esign/internal/config"
)

type Client struct {
	apiKey     string
	fromEmail  string
	templateID string
	testMode   bool
	httpClient *http.Client
}

type EmailPayload struct {
	ToEmail       string
	ToName        string
	Subject       string
	InvoiceNumber string
	InvoiceDate   string
	DueDate       string
	TotalAmount   string
	AttachmentPath string // path to HTML file to attach
}

func New(cfg *config.Config) *Client {
	return &Client{
		apiKey:     cfg.SendGridAPIKey,
		fromEmail:  cfg.FromEmail,
		templateID: cfg.SendGridTemplID,
		testMode:   cfg.TestMode,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Send dispatches an email with the invoice attached.
func (c *Client) Send(p EmailPayload) error {
	if c.testMode {
		fmt.Printf("  [TEST MODE] Would send to: %s (%s)\n", p.ToEmail, p.ToName)
		fmt.Printf("  [TEST MODE] Subject: %s\n", p.Subject)
		fmt.Printf("  [TEST MODE] Attachment: %s\n", p.AttachmentPath)
		return nil
	}

	if c.apiKey == "" {
		return fmt.Errorf("SENDGRID_API_KEY not configured")
	}

	// Build SendGrid v3 API request
	body, err := c.buildRequest(p)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("SendGrid API call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	respBody, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("SendGrid error (HTTP %d): %s", resp.StatusCode, string(respBody))
}

func (c *Client) buildRequest(p EmailPayload) ([]byte, error) {
	msg := map[string]any{
		"personalizations": []map[string]any{
			{
				"to": []map[string]string{
					{"email": p.ToEmail, "name": p.ToName},
				},
				"dynamic_template_data": map[string]string{
					"partner_name":  p.ToName,
					"invoice_no":    p.InvoiceNumber,
					"invoice_date":  p.InvoiceDate,
					"due_date":      p.DueDate,
					"total_amount":  p.TotalAmount,
				},
			},
		},
		"from": map[string]string{
			"email": c.fromEmail,
			"name":  "Invoice System",
		},
		"subject": p.Subject,
	}

	if c.templateID != "" {
		msg["template_id"] = c.templateID
	}

	// Attach the invoice file
	if p.AttachmentPath != "" {
		data, err := os.ReadFile(p.AttachmentPath)
		if err != nil {
			return nil, fmt.Errorf("reading attachment %s: %w", p.AttachmentPath, err)
		}
		encoded := base64.StdEncoding.EncodeToString(data)

		ext := filepath.Ext(p.AttachmentPath)
		mimeType := "text/html"
		if ext == ".pdf" {
			mimeType = "application/pdf"
		}

		msg["attachments"] = []map[string]string{
			{
				"content":     encoded,
				"filename":    filepath.Base(p.AttachmentPath),
				"type":        mimeType,
				"disposition": "attachment",
			},
		}
	}

	return json.Marshal(msg)
}
