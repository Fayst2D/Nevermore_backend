package email

import (
	"fmt"
	"net/smtp"
	"nevermore/config"
	"nevermore/pkg/logger"
)

type Service interface {
	SendVerificationEmail(to, token string) error
}

type service struct {
	cfg     config.SMTPConfig
	baseURL string
}

func New(cfg config.SMTPConfig, baseURL string) Service {
	return &service{
		cfg:     cfg,
		baseURL: baseURL,
	}
}

func (s *service) SendVerificationEmail(to, token string) error {
	from := s.cfg.From
	subject := "Подтверждение email адреса"
	body := fmt.Sprintf(`
Здравствуйте!

Для подтверждения вашего email адреса, пожалуйста, перейдите по следующей ссылке:

%s/auth/verify-email?token=%s

Если вы не регистрировались на нашем сайте, просто проигнорируйте это письмо.

С уважением,
Команда Nevermore
`, s.baseURL, token)

	msg := []byte(fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body)

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	log := logger.Get()
	log.Info().Msgf("Sending verification email to %s", to)

	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to send verification email to %s", to)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().Msgf("Verification email sent successfully to %s", to)
	return nil
}
