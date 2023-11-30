package xmail

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
)

// SendMail 发送邮件
// server := "smtp.example.tld:465"
// from := "username@example.tld"
// to := "username@anotherexample.tld"
func SendMail(server, from, password, to, subject, body string) error {
	// golang 发送邮件代码
	f := mail.Address{Address: from}
	t := mail.Address{Address: to}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = f.String()
	headers["To"] = t.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server

	host, _, _ := net.SplitHostPort(server)

	auth := smtp.PlainAuth("", from, password, host)

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", server, tlsConfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(f.Address); err != nil {
		return err
	}

	if err = c.Rcpt(t.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	if err != nil {
		log.Printf("Quit error: %v", err)
		return err
	}
	return nil
}
