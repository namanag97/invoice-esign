# invoice-esign

CLI tool that generates GST tax invoices from a partner CSV file and emails them to vendors via SendGrid.

Built for Indian GST compliance — handles CGST+SGST (intra-state) vs IGST (inter-state) automatically based on GSTIN state codes.

## Pipeline

```
Partner CSV
    │
    ▼
┌────────────┐
│   Parser   │  Read CSV, extract vendor records
└─────┬──────┘
      │
      ▼
┌────────────┐
│ Validator  │  GSTIN format, PAN, phone, email, amounts
└─────┬──────┘
      │
      ▼
┌────────────┐
│ Generator  │  Fill HTML template, compute GST, format INR
└─────┬──────┘
      │
      ▼
┌────────────┐
│   Mailer   │  SendGrid API with invoice attachment
└────────────┘
```

## Install

```bash
go install github.com/namanag97/invoice-esign/cmd/invoice-esign@latest
```

Or build from source:

```bash
git clone https://github.com/namanag97/invoice-esign.git
cd invoice-esign
go build -o invoice-esign ./cmd/invoice-esign
```

## Usage

### Full pipeline (generate + email)

```bash
invoice-esign run data/partner_data.csv
```

### Validate CSV only

```bash
invoice-esign validate data/partner_data.csv
```

### Generate invoices without emailing

```bash
invoice-esign generate data/partner_data.csv
```

### Options

```bash
invoice-esign run data/partners.csv --output ./results --template templates/custom.html
```

## CSV Format

| Column | Required | Description |
|--------|----------|-------------|
| `partner_name` | Yes | Vendor name |
| `partner_email` | Yes | Email to send invoice to |
| `gst_no` | Yes | 15-char Indian GSTIN |
| `payout_amount` | Yes | Taxable amount (supports commas) |
| `pan_no` | No | PAN number |
| `partner_phone` | No | Phone (auto-cleans +91 prefix) |
| `bank_account_number` | No | Bank account (masked in invoice) |
| `bank_accountifsccode` | No | IFSC code (used to identify bank name) |
| `partner_code` | No | Used in invoice number generation |
| `account_id` | No | Partner UUID |
| `utr` | No | Payment reference |

## GST Logic

- **Same state** (vendor GSTIN state code == company state code): CGST 9% + SGST 9%
- **Different state**: IGST 18%
- State is derived from the first 2 digits of the GSTIN
- Company default: Karnataka (29)

## Configuration

Set via environment variables or `invoice-esign.json`:

| Variable | Description | Default |
|----------|-------------|---------|
| `SENDGRID_API_KEY` | SendGrid API key | (required for real emails) |
| `SENDGRID_TEMPLATE_ID` | SendGrid dynamic template ID | - |
| `INVOICE_FROM_EMAIL` | Sender email address | - |
| `INVOICE_OUTPUT_DIR` | Output directory | `./output` |
| `INVOICE_TEST_MODE` | Set to `false` to send real emails | `true` |

## Output

```
output/
└── invoices/
    ├── VM-ZCUXXM-0126-001.html
    ├── VM-DCL6DX-0126-002.html
    └── ...
```

Each invoice is a self-contained HTML file with proper Indian tax invoice formatting, including GSTIN, HSN code (997156), bank details, amount in words, and signature block.

## License

MIT
