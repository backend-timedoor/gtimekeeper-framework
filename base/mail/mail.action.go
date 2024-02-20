package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"net/textproto"
	"reflect"

	"github.com/backend-timedoor/gtimekeeper-framework/base/job"
	"github.com/backend-timedoor/gtimekeeper-framework/base/job/custom"
	"github.com/backend-timedoor/gtimekeeper-framework/container"
	"github.com/jordan-wright/email"
)

var (
	mailStructAttributes []string
)

type Envelope interface {
	Content(any) Content
}

type Email struct {
	Config  *Config
	c       Envelope
	content Content
	data    any
}

type SendTo struct {
	To  []string
	Bcc []string
	Cc  []string
}

type Content struct {
	ReplyTo     []string
	Subject     string
	Text        []byte
	HTML        []byte
	Sender      string
	Headers     textproto.MIMEHeader
	ReadReceipt []string
}

func (e *Email) SendWithQueue(data []byte) error {
	email := email.Email{}
	errChan := make(chan error, 1)

	err := json.Unmarshal(data, &email)
	if err != nil {
		return fmt.Errorf("error unmarshal email data: %v", err)
	}

	go func(errChan chan error) {
		sender, auth := e.auth()
		if err = email.Send(sender, auth); err != nil {
			errChan <- fmt.Errorf("email send error: %v", err)
		}

		close(errChan)

	}(errChan)

	return <-errChan
}

func (e *Email) Send(c Envelope, data any) error {
	errChan := make(chan error, 1)
	e.c = c
	e.data = data
	e.content = c.Content(data)

	// struct validation
	err := e.mailStructValidation()
	if err != nil {
		return err
	}

	refAttr := reflect.ValueOf(c).Elem()
	sendTo := refAttr.FieldByName(mailStructAttributes[0]).Interface().(SendTo)

	html, err := e.templateBind()
	if err != nil {
		return fmt.Errorf("error create html template error: %v", err)
	}

	m := &email.Email{
		To:          sendTo.To,
		Bcc:         sendTo.Bcc,
		Cc:          sendTo.Cc,
		ReplyTo:     e.content.ReplyTo,
		From:        e.from(),
		Subject:     e.content.Subject,
		Text:        e.content.Text,
		HTML:        html,
		Sender:      e.content.Sender,
		Headers:     e.content.Headers,
		ReadReceipt: e.content.ReadReceipt,
	}

	go func(errChan chan error) {
		withQueueAttr := refAttr.FieldByName(mailStructAttributes[1]).Bool()
		if withQueueAttr {
			jobQueue := reflect.ValueOf(container.Get(job.ContainerName)).Interface().(*job.Job)
			err = jobQueue.Queue(&custom.EmailJob{}, m)
			if err != nil {
				errChan <- fmt.Errorf("error queue email: %v", err)
			}
		} else {
			sender, auth := e.auth()
			if err = m.Send(sender, auth); err != nil {
				errChan <- fmt.Errorf("email send error: %v", err)
			}
		}

		close(errChan)
	}(errChan)

	return <-errChan
}

func (e *Email) from() string {
	ref := reflect.TypeOf(e.c)
	_, ok := ref.MethodByName("View")

	if ok {
		from := reflect.ValueOf(e.c).MethodByName("From").Call([]reflect.Value{})

		return from[0].Interface().(string)
	}

	return e.Config.From
}

func (e *Email) templateBind() (html []byte, err error) {
	html = e.content.HTML
	ref := reflect.TypeOf(e.c)
	_, ok := ref.MethodByName("View")

	if ok {
		files := reflect.ValueOf(e.c).MethodByName("View").Call([]reflect.Value{})

		file := fmt.Sprintf("%s/%s/%s",
			e.Config.RootPath,
			e.Config.TemplatePath,
			files[0].Interface().(string),
		)

		template, err := template.ParseFiles(file)
		if err != nil {
			log.Fatal(err)
		}

		buf := new(bytes.Buffer)
		if err = template.Execute(buf, e.data); err != nil {
			log.Fatal(err)
		}

		html = buf.Bytes()
	}

	return
}

func (e *Email) auth() (string, smtp.Auth) {
	return fmt.Sprintf("%s:%d", e.Config.Host, e.Config.Port), smtp.PlainAuth(
		"",
		e.Config.Username,
		e.Config.Password,
		e.Config.Host,
	)
}

func (e *Email) mailStructValidation() error {
	el := reflect.ValueOf(e.c).Elem()
	fieldError := make([]string, 0)

	for _, attr := range mailStructAttributes {
		if !el.FieldByName(attr).IsValid() {
			fieldError = append(fieldError, attr)
		}
	}

	if len(fieldError) > 0 {
		return fmt.Errorf("mail struct must be containt attributes: %v", fieldError)
	}

	return nil
}
