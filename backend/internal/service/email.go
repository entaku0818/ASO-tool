package service

import (
	"fmt"
	"log"
	"os"

	resend "github.com/resend/resend-go/v2"
)

type EmailService struct {
	client  *resend.Client
	from    string
	enabled bool
}

func NewEmailService() *EmailService {
	apiKey := os.Getenv("RESEND_API_KEY")
	from := os.Getenv("RESEND_FROM_EMAIL")
	if from == "" {
		from = "noreply@aso-tool.app"
	}
	if apiKey == "" {
		log.Println("warn: RESEND_API_KEY is not set — email sending will be skipped")
		return &EmailService{enabled: false, from: from}
	}
	return &EmailService{
		client:  resend.NewClient(apiKey),
		from:    from,
		enabled: true,
	}
}

func (s *EmailService) SendLicenseKey(to, licenseKey string) error {
	if !s.enabled {
		log.Printf("info: [email-stub] to=%s license=%s", to, licenseKey)
		return nil
	}

	subject := "【ASO-tool macOS】ライセンスキーのお知らせ"
	html := fmt.Sprintf(`
<div style="font-family:sans-serif;max-width:560px;margin:0 auto">
  <h2>ご購入ありがとうございます</h2>
  <p>ASO-tool macOS版のライセンスキーをお届けします。</p>
  <div style="background:#f4f4f4;border-radius:8px;padding:20px;text-align:center;margin:24px 0">
    <code style="font-size:22px;font-weight:bold;letter-spacing:2px">%s</code>
  </div>
  <p>アプリを起動してキーとメールアドレスを入力するとアクティベートできます。</p>
  <p style="color:#888;font-size:12px">このキーは1つのアカウントにのみ有効です。</p>
</div>
`, licenseKey)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: subject,
		Html:    html,
	}
	_, err := s.client.Emails.Send(params)
	return err
}
