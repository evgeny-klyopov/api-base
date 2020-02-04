package service

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"
)

type MailHeaders = map[string]string

type authMail struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

type Sendmail struct {
	from    mail.Address
	to      mail.Address
	subject string
	body    string
	auth    map[string]authMail
}

func NewSendmail() (Sendmail, error) {
	var s Sendmail
	err := s.prepare()
	return s, err
}

func (s *Sendmail) prepare() error {
	config := []byte(os.Getenv("MAIL_CONFIG"))

	if err := json.Unmarshal(config, &s.auth); err != nil {
		return errors.New("bad mail config")
	}

	if len(s.auth) < 1 {
		return errors.New("not found setting mail")
	}

	return nil
}
func (s *Sendmail) SetFrom(name string, from string) *Sendmail {
	s.from = mail.Address{Name: name, Address: from}

	return s
}
func (s *Sendmail) SetTo(name string, to string) *Sendmail {
	s.to = mail.Address{Name: name, Address: to}

	return s
}
func (s *Sendmail) SetBody(body string) *Sendmail {
	s.body = body

	return s
}
func (s *Sendmail) SetSubject(subject string) *Sendmail {
	s.subject = subject

	return s
}

func (s *Sendmail) getMessages(headers MailHeaders) string {
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + s.body

	return message
}

func (s *Sendmail) getHeaders(addHeaders *MailHeaders) MailHeaders {
	var headers MailHeaders
	headers = make(map[string]string)
	headers["From"] = s.from.String()
	headers["To"] = s.to.String()
	headers["Subject"] = s.subject

	if addHeaders != nil {
		for key, value := range *addHeaders {
			headers[key] = value
		}
	}

	return headers
}
func (s *Sendmail) GetConfigByMailKey(mailConfigKey string) (*authMail, error) {
	config := s.auth[mailConfigKey]

	if (authMail{}) == config {
		return nil, fmt.Errorf("not found config for key - %s", mailConfigKey)
	}

	return &config, nil
}
func (s *Sendmail) Send(mailConfigKey string, addHeaders *MailHeaders) error {
	config, err := s.GetConfigByMailKey(mailConfigKey)
	if err != nil {
		return err
	}

	headers := s.getHeaders(addHeaders)
	message := s.getMessages(headers)
	servername := net.JoinHostPort(config.Host, config.Port)

	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	//TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         config.Host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsConfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(s.from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(s.to.Address); err != nil {
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
		return err
	}

	return nil
}
