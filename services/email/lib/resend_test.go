package email

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/resend/resend-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Resend client
type MockEmailsService struct {
	mock.Mock
}

func (m *MockEmailsService) Send(params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*resend.SendEmailResponse), args.Error(1)
}

func (m *MockEmailsService) SendWithContext(ctx context.Context, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*resend.SendEmailResponse), args.Error(1)
}

func (m *MockEmailsService) Get(emailId string) (*resend.Email, error) {
	args := m.Called(emailId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*resend.Email), args.Error(1)
}

func (m *MockEmailsService) GetWithContext(ctx context.Context, emailId string) (*resend.Email, error) {
	args := m.Called(ctx, emailId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*resend.Email), args.Error(1)
}

func (m *MockEmailsService) Cancel(emailId string) (*resend.CancelScheduledEmailResponse, error) {
	args := m.Called(emailId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*resend.CancelScheduledEmailResponse), args.Error(1)
}

func (m *MockEmailsService) CancelWithContext(ctx context.Context, emailId string) (*resend.CancelScheduledEmailResponse, error) {
	args := m.Called(ctx, emailId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*resend.CancelScheduledEmailResponse), args.Error(1)
}

func (m *MockEmailsService) List() (resend.ListEmailsResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return resend.ListEmailsResponse{}, args.Error(1)
	}
	return args.Get(0).(resend.ListEmailsResponse), args.Error(1)
}

func (m *MockEmailsService) ListWithContext(ctx context.Context) (resend.ListEmailsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return resend.ListEmailsResponse{}, args.Error(1)
	}
	return args.Get(0).(resend.ListEmailsResponse), args.Error(1)
}

func (m *MockEmailsService) ListWithOptions(ctx context.Context, options *resend.ListOptions) (resend.ListEmailsResponse, error) {
	args := m.Called(ctx, options)
	if args.Get(0) == nil {
		return resend.ListEmailsResponse{}, args.Error(1)
	}
	return args.Get(0).(resend.ListEmailsResponse), args.Error(1)
}

func (m *MockEmailsService) SendWithOptions(ctx context.Context, params *resend.SendEmailRequest, options *resend.SendEmailOptions) (*resend.SendEmailResponse, error) {
	args := m.Called(ctx, params, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*resend.SendEmailResponse), args.Error(1)
}

func (m *MockEmailsService) Update(params *resend.UpdateEmailRequest) (*resend.UpdateEmailResponse, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*resend.UpdateEmailResponse), args.Error(1)
}

func (m *MockEmailsService) UpdateWithContext(ctx context.Context, params *resend.UpdateEmailRequest) (*resend.UpdateEmailResponse, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*resend.UpdateEmailResponse), args.Error(1)
}

func TestResend(t *testing.T) {
	t.Run("should return error if RESEND_API_KEY is not set", func(t *testing.T) {
		os.Unsetenv("RESEND_API_KEY")
		adapter, err := Resend()
		assert.Error(t, err)
		assert.Nil(t, adapter)
		assert.Equal(t, "RESEND_API_KEY environment variable is not set", err.Error())
	})

	t.Run("should return error if RESEND_DEFAULT_FROM is not set", func(t *testing.T) {
		os.Setenv("RESEND_API_KEY", "test-key")
		defer os.Unsetenv("RESEND_API_KEY")
		os.Unsetenv("RESEND_DEFAULT_FROM")

		adapter, err := Resend()

		assert.Error(t, err)
		assert.Nil(t, adapter)
		assert.Equal(t, "RESEND_DEFAULT_FROM environment variable is not set", err.Error())
	})

	t.Run("should return a new ResendAdapter", func(t *testing.T) {
		os.Setenv("RESEND_API_KEY", "test-key")
		os.Setenv("RESEND_DEFAULT_FROM", "test@example.com")
		defer os.Unsetenv("RESEND_API_KEY")
		defer os.Unsetenv("RESEND_DEFAULT_FROM")

		adapter, err := Resend()

		assert.NoError(t, err)
		assert.NotNil(t, adapter)
		assert.Equal(t, "test@example.com", adapter.defaultFrom)
	})
}

func TestResendAdapter_Send(t *testing.T) {
	apiKey := "test-key"
	defaultFrom := "default@example.com"

	t.Run("should return error if no recipients", func(t *testing.T) {
		adapter := &ResendAdapter{
			client:      resend.NewClient(apiKey),
			defaultFrom: defaultFrom,
		}
		data := EmailData{
			To:       []string{},
			Subject:  "Test",
			Html: "<h1>Test</h1>",
		}

		err := adapter.Send(context.Background(), data)

		assert.Error(t, err)
		assert.Equal(t, "email must have at least one recipient", err.Error())
	})

	t.Run("should use default from address if not provided", func(t *testing.T) {
		mockEmailsService := new(MockEmailsService)
		client := &resend.Client{
			Emails: mockEmailsService,
		}
		adapter := &ResendAdapter{
			client:      client,
			defaultFrom: defaultFrom,
		}
		data := EmailData{
			To:       []string{"recipient@example.com"},
			Subject:  "Test",
			Html: "<h1>Test</h1>",
		}

		expectedParams := &resend.SendEmailRequest{
			From:    defaultFrom,
			To:      data.To,
			Subject: data.Subject,
			Html:    data.Html,
		}

		mockEmailsService.On("SendWithContext", context.Background(), expectedParams).Return(&resend.SendEmailResponse{Id: "test-id"}, nil)

		err := adapter.Send(context.Background(), data)

		assert.NoError(t, err)
		mockEmailsService.AssertExpectations(t)
	})

	t.Run("should send email successfully", func(t *testing.T) {
		mockEmailsService := new(MockEmailsService)
		client := &resend.Client{
			Emails: mockEmailsService,
		}
		adapter := &ResendAdapter{
			client:      client,
			defaultFrom: defaultFrom,
		}
		data := EmailData{
			From:     "custom@example.com",
			To:       []string{"recipient@example.com"},
			Subject:  "Test",
			Html: "<h1>Test</h1>",
		}

		expectedParams := &resend.SendEmailRequest{
			From:    data.From,
			To:      data.To,
			Subject: data.Subject,
			Html:    data.Html,
		}

		mockEmailsService.On("SendWithContext", context.Background(), expectedParams).Return(&resend.SendEmailResponse{Id: "test-id"}, nil)

		err := adapter.Send(context.Background(), data)

		assert.NoError(t, err)
		mockEmailsService.AssertExpectations(t)
	})

	t.Run("should return error on failed send", func(t *testing.T) {
		mockEmailsService := new(MockEmailsService)
		client := &resend.Client{
			Emails: mockEmailsService,
		}
		adapter := &ResendAdapter{
			client:      client,
			defaultFrom: defaultFrom,
		}
		data := EmailData{
			To:       []string{"recipient@example.com"},
			Subject:  "Test",
			Html: "<h1>Test</h1>",
		}

		expectedParams := &resend.SendEmailRequest{
			From:    defaultFrom,
			To:      data.To,
			Subject: data.Subject,
			Html:    data.Html,
		}
		sendErr := errors.New("send failed")

		mockEmailsService.On("SendWithContext", context.Background(), expectedParams).Return(nil, sendErr)

		err := adapter.Send(context.Background(), data)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to send email via resend")
		mockEmailsService.AssertExpectations(t)
	})
}

