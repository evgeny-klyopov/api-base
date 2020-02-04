package common

import (
	"api-base/service"
	"errors"
	"net/mail"
)

type mailSender struct {
	mailConfigKey string
	allowedKeys []string
	sender service.Sendmail
	defaultEmail mail.Address
}

const UseDefaultTo = "to"
const UseDefaultFrom = "from"

const SupportMailConfigKey = "support"
const InfoMailConfigKey = "info"

func setEmail(item *mail.Address, key string, defaultAddress *[]string, defaultEmail mail.Address) *mail.Address {
	var exists bool
	exists, _ = InArray(key, *defaultAddress)
	if exists == true {
		item = &defaultEmail
	}

	return item
}

func SendMail(mailConfigKey string, subject string, message string, defaultAddress *[]string, to *mail.Address, from *mail.Address, headers *service.MailHeaders) error {
	var obj mailSender

	obj.mailConfigKey = mailConfigKey
	obj.allowedKeys = []string{SupportMailConfigKey, InfoMailConfigKey}

	err := obj.checkKey()
	if err != nil {
		return err
	}

	if nil != defaultAddress {
		to = setEmail(to, UseDefaultTo, defaultAddress, obj.defaultEmail)
		from = setEmail(from, UseDefaultFrom, defaultAddress, obj.defaultEmail)
	}

	obj.sender.SetTo(to.Name, to.Address)
	obj.sender.SetFrom(from.Name, from.Address)
	obj.sender.SetSubject(subject)
	obj.sender.SetBody(message)

	err = obj.sender.Send(obj.mailConfigKey, headers)
	if err != nil {
		return err
	}

	return nil
}

func (m *mailSender) checkKey() error {
	exists, _ := InArray(m.mailConfigKey, m.allowedKeys)

	if false == exists {
		return errors.New("not found config mail")
	}

	err := m.prepare()
	if err != nil {
		return err
	}

	return nil
}
func (m *mailSender) prepare() error {
	var err error
	m.sender, err = service.NewSendmail()
	if err != nil {
		return err
	}

	config, err := m.sender.GetConfigByMailKey(m.mailConfigKey)
	if err != nil {
		return err
	}

	m.defaultEmail = mail.Address{
		Name:    config.Name,
		Address: config.Username,
	}

	return nil
}