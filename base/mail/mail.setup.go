package mail

import "github.com/backend-timedoor/gtimekeeper-framework/base/contracts"

func BootMail() contracts.Mail {
	return &Email{}
}