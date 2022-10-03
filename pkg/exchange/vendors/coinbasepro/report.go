package coinbasepro

import (
	"context"
	"errors"
	"fmt"
)

type ReportSpec struct {
	AccountID string       `json:"account_id"`
	EndDate   Time         `json:"end_date"`
	Email     string       `json:"email"`
	Format    ReportFormat `json:"format"`
	ProductID ProductID    `json:"product_id"`
	StartDate Time         `json:"start_date"`
	Type      ReportType   `json:"type"`
}

func (c *ReportSpec) Validate() error {
	if c.Type != ReportTypeAccount && c.Type != ReportTypeFills {
		return errors.New("report 'type' must be one of 'account' or 'fills")
	}
	if c.Type == ReportTypeFills && c.ProductID == "" {
		return errors.New("'product_id' required for report type 'fills'")
	}
	if c.Type == ReportTypeAccount && c.AccountID == "" {
		return errors.New("'account_id' required for report type 'account'")
	}
	if c.Format != ReportFormatPDF && c.Format != ReportFormatCSV {
		return errors.New("'format' must be one of 'pdf' or 'csv'")
	}
	if c.EndDate.Time().IsZero() {
		return errors.New("'end_date' is required")
	}
	if c.StartDate.Time().IsZero() {
		return errors.New("'start_date' is required")
	}
	return nil
}

type ReportFormat string

const (
	ReportFormatCSV ReportFormat = "csv"
	ReportFormatPDF ReportFormat = "pdf" // default
)

type ReportType string

const (
	ReportTypeAccount ReportType = "account"
	ReportTypeFills   ReportType = "fills"
)

type Report struct {
	CompletedAt Time            `json:"completed_at"`
	CreatedAt   Time            `json:"created_at"`
	ExpiresAt   Time            `json:"expires_at"`
	FileURL     string          `json:"file_url"`
	ID          string          `json:"id"`
	Params      ReportDateRange `json:"params"`
	Status      ReportStatus    `json:"status"`
	Type        ReportType      `json:"type"`
}

type ReportDateRange struct {
	EndDate   Time `json:"end_date"`
	StartDate Time `json:"start_date"`
}

type ReportStatus string

const (
	// ReportStatusCreating indicates that report is being created
	ReportStatusCreating ReportStatus = "creating"
	// ReportStatusPending indicates that the report request has been accepted and is awaiting processing
	ReportStatusPending ReportStatus = "pending"
	// ReportStatusReady indicates that the report is ready for download from `file_url`
	ReportStatusReady ReportStatus = "ready"
)

// CreateReport creates request for batches of historic Profile information in various human and machine readable forms.
// Reports will be generated when resources are available. Report status can be queried using GetReport.
func (c *CoinbasePro) CreateReport(ctx context.Context, createReportSpec ReportSpec) (Report, error) {
	result := struct {
		ID string `json:"id"`
	}{}
	path := fmt.Sprintf("/%s/", coinbaseproReports)
	err := c.API.Post(ctx, path, createReportSpec, &result)
	if err != nil {
		return Report{}, err
	}
	// POST coinbasepro response is partial; retrieve full representation
	return c.GetReport(ctx, result.ID)
}

// GetReport retrieves the status of the processing of a Report request. When the ReportStatus is 'ready',
// the Report will be available for download at the FileURL.
func (c *CoinbasePro) GetReport(ctx context.Context, reportID string) (Report, error) {
	var report Report
	path := fmt.Sprintf("/%s/%s", coinbaseproReports, reportID)
	return report, c.API.Get(ctx, path, &report)
}
