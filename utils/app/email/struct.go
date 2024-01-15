package email

import "net/textproto"

type SendTo struct{
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
