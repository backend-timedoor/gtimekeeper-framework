package mail

import "github.com/backend-timedoor/gtimekeeper-framework/container"

const ContainerName string = "mail"

type Config struct {
	Host         string
	Port         int
	Username     string
	Password     string
	From         string
	RootPath     string
	TemplatePath string
}

func New(config *Config) *Email {
	mailStructAttributes = []string{"SendTo", "WithQueue", "Attachments"}
	e := &Email{
		Config: config,
	}

	container.Set(ContainerName, e)

	return e
}
