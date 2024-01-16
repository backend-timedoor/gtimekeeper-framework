package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"reflect"

	"github.com/backend-timedoor/gtimekeeper-framework/app"
	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	envelop "github.com/backend-timedoor/gtimekeeper-framework/utils/app/email"
	"github.com/jordan-wright/email"
)

type Email struct {
	c contracts.Email
	content envelop.Content
	data any
}

func (e Email) Send(c contracts.Email, data any) {
	e.c = c
	e.data = data
	e.content = c.Content(data)

	r := reflect.ValueOf(c).Elem()
	sendTo := r.Field(0).Interface().(envelop.SendTo)

	em := &email.Email{
		To: sendTo.To,
		Bcc: sendTo.Bcc,
		Cc: sendTo.Cc,
		ReplyTo: e.content.ReplyTo,
		From: "Edwin Diradinata <edwindiradinata@gmail.com>",
		Subject: e.content.Subject,
		Text: e.content.Text,
		HTML:  e.TemplateBind(),
		Sender: e.content.Sender,
		Headers: e.content.Headers,
		ReadReceipt: e.content.ReadReceipt,
	}
	
	host := app.Config.GetString("mail.host")
	sender := fmt.Sprintf("%s:%s", host, app.Config.GetString("mail.port"))
	auth := smtp.PlainAuth(
		"",
		app.Config.GetString("mail.username"),
		app.Config.GetString("mail.password"),
		host,
	)

	em.Send(sender, auth)
}

func (e Email) TemplateBind() (html []byte) {
	html = e.content.HTML
	ref := reflect.TypeOf(e.c)
	_, ok := ref.MethodByName("View")

	if ok {
		files := reflect.ValueOf(e.c).MethodByName("View").Call([]reflect.Value{})

		file :=  fmt.Sprintf("%s/%s/%s", 
			app.Config.GetString("path.root"),
			app.Config.GetString("path.mail"),
			files[0].Interface().(string),
		)

		template, err := template.ParseFiles(file)
		if err != nil {
			app.Log.Error(err)
		}

		buf := new(bytes.Buffer)
		if err = template.Execute(buf, e.data); err != nil {
			app.Log.Error(err)
		}

		html = buf.Bytes()
	}

	return
}