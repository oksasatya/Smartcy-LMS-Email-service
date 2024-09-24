package service

import (
	"context"
	"fmt"
	"github.com/ghssni/Smartcy-LMS/Email-Service/internal/models"
	"github.com/ghssni/Smartcy-LMS/Email-Service/internal/repository"
	pb "github.com/ghssni/Smartcy-LMS/Email-Service/pb/proto"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"os"
	"time"
)

type EmailService struct {
	pb.UnimplementedEmailServiceServer
	emailRepo    repository.EmailsRepository
	emailLogRepo repository.EmailsLogRepository
}

func NewEmailService(emailRepo repository.EmailsRepository, emailLogRepo repository.EmailsLogRepository) *EmailService {
	return &EmailService{
		emailRepo:    emailRepo,
		emailLogRepo: emailLogRepo,
	}
}

func (s *EmailService) SendPaymentDueEmailRequest(ctx context.Context, req *pb.SendPaymentDueEmailRequest) (*pb.SendPaymentDueEmailResponse, error) {
	email := &models.Email{
		UserID:    req.UserId,
		EmailType: "payment_due",
		Email:     req.Email,
	}

	_, err := s.emailRepo.InsertEmail(email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert email: %v", err)
	}

	// send email
	err = SendEmailPayment(req.Email, req.CourseName, req.PaymentLink)
	statusStr := "sent"
	errorMsg := ""
	if err != nil {
		logrus.Println("Error sending email:", err)
		statusStr = "failed"
		errorMsg = err.Error()
	}

	// log email
	emailLog := &models.EmailLogs{
		UserID:       req.UserId,
		Email:        req.Email,
		Status:       statusStr,
		SentAt:       time.Now(),
		ErrorMessage: errorMsg,
	}

	_, err = s.emailLogRepo.InsertEmailLog(emailLog)

	response := &pb.SendPaymentDueEmailResponse{
		Meta: &pb.Meta{
			Code:    int32(codes.OK),
			Message: "Email sent successfully",
			Status:  http.StatusText(http.StatusOK),
		},
		Success: statusStr == "sent",
	}
	return response, nil
}

func (s *EmailService) SendForgotPasswordEmailRequest(ctx context.Context, req *pb.SendForgotPasswordEmailRequest) (*pb.SendForgotPasswordEmailResponse, error) {
	email := &models.Email{
		UserID:    req.UserId,
		EmailType: "forgot_password",
		Email:     req.Email,
	}

	_, err := s.emailRepo.InsertEmail(email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert email: %v", err)
	}

	// send email
	err = SendEmailForgotPassword(req.Email, req.ResetLink, req.ResetToken)
	statusStr := "sent"
	errorMsg := ""
	if err != nil {
		logrus.Println("Error sending email:", err)
		statusStr = "failed"
		errorMsg = err.Error()
	}

	// log email
	emailLog := &models.EmailLogs{
		UserID:       req.UserId,
		Email:        req.Email,
		Status:       statusStr,
		SentAt:       time.Now(),
		ErrorMessage: errorMsg,
	}

	_, err = s.emailLogRepo.InsertEmailLog(emailLog)

	response := &pb.SendForgotPasswordEmailResponse{
		Meta: &pb.Meta{
			Code:    int32(codes.OK),
			Message: "Email sent successfully",
			Status:  http.StatusText(http.StatusOK),
		},
		Success: statusStr == "sent",
	}

	return response, nil
}

func (s *EmailService) SendPaymentSuccessEmailRequest(ctx context.Context, req *pb.SendPaymentSuccessEmailRequest) (*pb.SendPaymentSuccessEmailResponse, error) {
	email := &models.Email{
		UserID:    req.UserId,
		EmailType: "payment_success",
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	_, err := s.emailRepo.InsertEmail(email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert email: %v", err)
	}

	// send email
	err = SendEmailSuccess(req.Email, req.CourseName)
	statusStr := "sent"
	errorMsg := ""
	if err != nil {
		logrus.Println("Error sending email:", err)
		statusStr = "failed"
		errorMsg = err.Error()
	}

	// log email
	emailLog := &models.EmailLogs{
		UserID:       req.UserId,
		Email:        req.Email,
		Status:       statusStr,
		SentAt:       time.Now(),
		ErrorMessage: errorMsg,
	}

	_, err = s.emailLogRepo.InsertEmailLog(emailLog)

	response := &pb.SendPaymentSuccessEmailResponse{
		Meta: &pb.Meta{
			Code:    int32(codes.OK),
			Message: "Email sent successfully",
			Status:  http.StatusText(http.StatusOK),
		},
		Success: statusStr == "sent",
	}

	return response, nil
}

// SendEmailPayment sends an email to the student with a payment URL
func SendEmailPayment(email, courseName, paymentURL string) error {
	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAILGUN_API_KEY")

	mg := mailgun.NewMailgun(domain, apiKey)

	sender := fmt.Sprintf("Smartcy LMS <no-reply@%s>", domain)
	subject := "Payment Confirmation for Course"
	htmlBody := fmt.Sprintf(`
		<html>
			<body>
				<h2>Payment Confirmation</h2>
				<p>Dear Student,</p>
				<p>Your payment for the course <b>%s</b> is pending.</p>
				<p>Please complete the payment using the following link:</p>
				<a href="%s">Click here to pay</a>
				<p>Thank you!</p>
			</body>
		</html>`, courseName, paymentURL)
	recipient := email

	message := mg.NewMessage(sender, subject, "", recipient)

	message.SetHtml(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	if err != nil {
		logrus.Println("Error sending email:", err)
		return err
	}
	return nil
}

// SendEmailSuccess sends an email to the student confirming the payment
func SendEmailSuccess(email, courseName string) error {
	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAILGUN_API_KEY")

	mg := mailgun.NewMailgun(domain, apiKey)

	sender := fmt.Sprintf("Smartcy LMS <no-reply@%s>", domain)
	subject := "Payment Confirmation for Course"
	htmlBody := fmt.Sprintf(`
		<html>
			<body>
				<h2>Payment Confirmation</h2>
				<p>Dear Student,</p>
				<p>Your payment for the course <b>%s</b> has been successfully processed.</p>
				<p>Thank you for your payment!</p>
			</body>
		</html>`, courseName)
	recipient := email

	message := mg.NewMessage(sender, subject, "", recipient)

	message.SetHtml(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	if err != nil {
		logrus.Println("Error sending email:", err)
		return err
	}
	return nil
}

// SendEmailForgotPassword sends an email to the student with a password reset link
func SendEmailForgotPassword(email, resetURL, resetToken string) error {
	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAILGUN_API_KEY")

	mg := mailgun.NewMailgun(domain, apiKey)
	fullResetURL := fmt.Sprintf("%s?token=%s", resetURL, resetToken)

	sender := fmt.Sprintf("Smartcy LMS <no-reply@%s>", domain)
	subject := "Reset Password Request"
	htmlBody := fmt.Sprintf(`
		<html>
			<body>
				<h2>Reset Password</h2>
				<p>Dear Student,</p>
				<p>We have received a request to reset your password.</p>
				<p>Please click the link below to reset your password:</p>
				<a href="%s">Reset Password</a>
				<p>If you did not request this, please ignore this email.</p>
				<p>Thank you!</p>
			</body>
		</html>`, fullResetURL)
	recipient := email

	message := mg.NewMessage(sender, subject, "", recipient)

	message.SetHtml(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	if err != nil {
		logrus.Println("Error sending email:", err)
		return err
	}
	return nil
}
