package email

import (
	"fmt"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

const (
	welcomeSubject = "Welcome to JIJI!"
)

const welcomeText = `
	Hi there!

	Welcome to JIJI! We really hope you enjoy using our application!

	Best,
	Yuichi
`

const welcomeHTML = `
	Hi there!<br/>
	<br/>
	Welcome to <a href="https://www.jiji_demo.com">JIJI</a>! We really hope you enjoy using our application!<br/>
	<br/>
	Best,<br/>
	Yuichi
`

func WithMailgun(domain, apiKey, publicKey string) ClientConfig {
	return func(client *Client) {
		mg := mailgun.NewMailgun(domain, apiKey, publicKey)
		client.mg = mg
	}
}

func WithSender(username, email string) ClientConfig {
	return func(client *Client) {
		client.from = buildEmail(username, email)
	}
}

type ClientConfig func(*Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		// Set a default from email address...
		from: "support@JIJI",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

type Client struct {
	from string
	mg   mailgun.Mailgun
}

func (client *Client) Welcome(toUsername, toEmail string) error {
	message := mailgun.NewMessage(client.from, welcomeSubject, welcomeText, buildEmail(toUsername, toEmail))
	message.SetHtml(welcomeHTML)
	_, _, err := client.mg.Send(message)
	// if err != nil {
	// 	panic(err)
	// }
	return err
}

func buildEmail(username, email string) string {
	if username == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", username, email)
}
