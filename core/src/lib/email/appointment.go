package email

import (
	"context"
	"fmt"
	"log"
	"mynute-go/core/src/config/db/model"

	"gorm.io/gorm"
)

// AppointmentEmailData contains all data needed for appointment emails
type AppointmentEmailData struct {
	ClientName      string
	ClientEmail     string
	EmployeeName    string
	EmployeeEmail   string
	ServiceName     string
	AppointmentDate string
	AppointmentTime string
	Duration        string
	BranchAddress   string
	Language        string
}

// AppointmentEmailService handles sending appointment-related emails
type AppointmentEmailService struct {
	sender           Sender
	templateRenderer *TemplateRenderer
}

// NewAppointmentEmailService creates a new appointment email service
func NewAppointmentEmailService(sender Sender, templateDir, translationDir string) *AppointmentEmailService {
	return &AppointmentEmailService{
		sender:           sender,
		templateRenderer: NewTemplateRenderer(templateDir, translationDir),
	}
}

// LoadAppointmentData loads all necessary data for an appointment email
func (s *AppointmentEmailService) LoadAppointmentData(tx *gorm.DB, appointment *model.Appointment, language string) (*AppointmentEmailData, error) {
	// Load client
	var client model.Client
	if err := tx.Model(&model.Client{}).Where("id = ?", appointment.ClientID).First(&client).Error; err != nil {
		return nil, fmt.Errorf("failed to load client: %w", err)
	}

	// Load employee
	var employee model.Employee
	if err := tx.Model(&model.Employee{}).Where("id = ?", appointment.EmployeeID).First(&employee).Error; err != nil {
		return nil, fmt.Errorf("failed to load employee: %w", err)
	}

	// Load service
	var service model.Service
	if err := tx.Model(&model.Service{}).Where("id = ?", appointment.ServiceID).First(&service).Error; err != nil {
		return nil, fmt.Errorf("failed to load service: %w", err)
	}

	// Load branch
	var branch model.Branch
	if err := tx.Model(&model.Branch{}).Where("id = ?", appointment.BranchID).First(&branch).Error; err != nil {
		return nil, fmt.Errorf("failed to load branch: %w", err)
	}

	// Format times
	startTime := appointment.StartTime
	endTime := appointment.EndTime
	duration := endTime.Sub(startTime)

	// Use provided language (defaults to "en" if empty)
	if language == "" {
		language = "en"
	}

	// Build branch address
	branchAddress := branch.Street
	if branch.Number != "" {
		branchAddress = fmt.Sprintf("%s, %s", branchAddress, branch.Number)
	}
	if branch.Neighborhood != "" {
		branchAddress = fmt.Sprintf("%s, %s", branchAddress, branch.Neighborhood)
	}
	if branch.City != "" {
		branchAddress = fmt.Sprintf("%s - %s", branchAddress, branch.City)
	}
	if branch.State != "" {
		branchAddress = fmt.Sprintf("%s, %s", branchAddress, branch.State)
	}

	return &AppointmentEmailData{
		ClientName:      fmt.Sprintf("%s %s", client.Name, client.Surname),
		ClientEmail:     client.Email,
		EmployeeName:    fmt.Sprintf("%s %s", employee.Name, employee.Surname),
		EmployeeEmail:   employee.Email,
		ServiceName:     service.Name,
		AppointmentDate: startTime.Format("Monday, January 2, 2006"),
		AppointmentTime: fmt.Sprintf("%s - %s", startTime.Format("3:04 PM"), endTime.Format("3:04 PM")),
		Duration:        fmt.Sprintf("%d minutes", int(duration.Minutes())),
		BranchAddress:   branchAddress,
		Language:        language,
	}, nil
}

// SendAppointmentCreatedEmails sends creation emails to both client and employee
func (s *AppointmentEmailService) SendAppointmentCreatedEmails(ctx context.Context, tx *gorm.DB, appointment *model.Appointment, language string) error {
	data, err := s.LoadAppointmentData(tx, appointment, language)
	if err != nil {
		return fmt.Errorf("failed to load appointment data: %w", err)
	}

	// Send email to client
	if err := s.sendEmail(ctx, "appointment_created", data.ClientEmail, data.ClientName, data); err != nil {
		log.Printf("Failed to send appointment created email to client %s: %v", data.ClientEmail, err)
		// Don't return error - continue to send employee email
	}

	// Send email to employee
	if err := s.sendEmail(ctx, "appointment_created", data.EmployeeEmail, data.EmployeeName, data); err != nil {
		log.Printf("Failed to send appointment created email to employee %s: %v", data.EmployeeEmail, err)
	}

	return nil
}

// SendAppointmentUpdatedEmails sends update emails to both client and employee
func (s *AppointmentEmailService) SendAppointmentUpdatedEmails(ctx context.Context, tx *gorm.DB, appointment *model.Appointment, language string) error {
	data, err := s.LoadAppointmentData(tx, appointment, language)
	if err != nil {
		return fmt.Errorf("failed to load appointment data: %w", err)
	}

	// Send email to client
	if err := s.sendEmail(ctx, "appointment_updated", data.ClientEmail, data.ClientName, data); err != nil {
		log.Printf("Failed to send appointment updated email to client %s: %v", data.ClientEmail, err)
	}

	// Send email to employee
	if err := s.sendEmail(ctx, "appointment_updated", data.EmployeeEmail, data.EmployeeName, data); err != nil {
		log.Printf("Failed to send appointment updated email to employee %s: %v", data.EmployeeEmail, err)
	}

	return nil
}

// SendAppointmentCancelledEmails sends cancellation emails to both client and employee
func (s *AppointmentEmailService) SendAppointmentCancelledEmails(ctx context.Context, tx *gorm.DB, appointment *model.Appointment, language string) error {
	data, err := s.LoadAppointmentData(tx, appointment, language)
	if err != nil {
		return fmt.Errorf("failed to load appointment data: %w", err)
	}

	// Send email to client
	if err := s.sendEmail(ctx, "appointment_cancelled", data.ClientEmail, data.ClientName, data); err != nil {
		log.Printf("Failed to send appointment cancelled email to client %s: %v", data.ClientEmail, err)
	}

	// Send email to employee
	if err := s.sendEmail(ctx, "appointment_cancelled", data.EmployeeEmail, data.EmployeeName, data); err != nil {
		log.Printf("Failed to send appointment cancelled email to employee %s: %v", data.EmployeeEmail, err)
	}

	return nil
}

// sendEmail is a helper function to render and send an email
func (s *AppointmentEmailService) sendEmail(ctx context.Context, templateName, toEmail, recipientName string, data *AppointmentEmailData) error {
	// Create template data
	templateData := TemplateData{
		"ClientName":      data.ClientName,
		"ServiceName":     data.ServiceName,
		"EmployeeName":    data.EmployeeName,
		"AppointmentDate": data.AppointmentDate,
		"AppointmentTime": data.AppointmentTime,
		"Duration":        data.Duration,
		"BranchAddress":   data.BranchAddress,
	}

	// Render email
	rendered, err := s.templateRenderer.RenderEmail(templateName, data.Language, templateData)
	if err != nil {
		return fmt.Errorf("failed to render email: %w", err)
	}

	// Send email
	emailData := EmailData{
		To:      []string{toEmail},
		Subject: rendered.Subject,
		Html:    rendered.HTMLBody,
	}

	if err := s.sender.Send(ctx, emailData); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Successfully sent %s email to %s", templateName, toEmail)
	return nil
}
