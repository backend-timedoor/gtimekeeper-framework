package providers

import (
	"github.com/backend-timedoor/gtimekeeper-framework/app"
	"github.com/backend-timedoor/gtimekeeper-framework/base/mail"
)

type MailServiceProvider struct{}

func (p *MailServiceProvider) Boot() {}

func (p *MailServiceProvider) Register() {
	app.Mail = mail.BootMail()
}
