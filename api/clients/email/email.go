package email

import (
	"bytes"
	"context"
	_ "embed"
	"text/template"
)

//go:embed generic_email_template_html.tmpl
var htmlTemplate string

//go:embed generic_email_template_text.tmpl
var textTemplate string

type SendGenericEmailParams struct {
	To            string
	Subject       string
	Title         string
	BodyHTML      string
	BodyPlainText string
	CTAText       string
	CTALink       string
	Reason        string
}

func SendGenericEmail(ctx context.Context, params SendGenericEmailParams) error {
	htmlEmail, textEmail, err := BuildEmail(params)
	if err != nil {
		return err
	}

	SendEmailUsingSES(ctx, AWSEmailParams{
		Email:    params.To,
		Subject:  params.Subject,
		BodyHTML: htmlEmail,
		BodyText: textEmail,
	})

	return nil
}

var SenderDetails = SenderParams{
	Name:            "Abstractions",
	Link:            "https://tryabstractions.com",
	Logo:            "https://tryabstractions.com/assets/logo.svg",
	LegalName:       "Artery Software Labs",
	Address:         "Bangalore, India",
	UnsubscribeLink: "https://tryabstractions.com/unsubscribe",
	PreferencesLink: "https://tryabstractions.com/preferences",
	PrimaryColor:    "#228ff4",
	Email:           "support@tryabstractions.com",
}

type SenderParams struct {
	Name            string
	Link            string
	Logo            string
	LegalName       string
	Address         string
	UnsubscribeLink string
	PreferencesLink string
	PrimaryColor    string
	Email           string
}

type GenericEmailTemplateParams struct {
	Title           string
	BodyHTML        string
	BodyPlainText   string
	CTAText         string
	CTALink         string
	EmailReason     string
	SenderName      string
	SenderLink      string
	SenderLogo      string
	SenderLegalName string
	SenderAddress   string
	UnsubscribeLink string
	PreferencesLink string
	PrimaryColor    string
}

// BuildEmail takes the generic email parameters and returns the HTML and plain text email bodies
// It also returns an error if the email cannot be built
func BuildEmail(params SendGenericEmailParams) (string, string, error) {

	// read the generic_email_template.tmpl file and replace it with the variables
	htmlEmailTemplate, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}

	textEmailTemplate, err := template.New("email").Parse(textTemplate)
	if err != nil {
		panic(err)
	}

	templateParams := GenericEmailTemplateParams{
		Title:           params.Title,
		BodyHTML:        params.BodyHTML,
		BodyPlainText:   params.BodyPlainText,
		CTAText:         params.CTAText,
		CTALink:         params.CTALink,
		EmailReason:     params.Reason,
		SenderName:      SenderDetails.Name,
		SenderLink:      SenderDetails.Link,
		SenderLogo:      SenderDetails.Logo,
		SenderLegalName: SenderDetails.LegalName,
		SenderAddress:   SenderDetails.Address,
		UnsubscribeLink: SenderDetails.UnsubscribeLink,
		PreferencesLink: SenderDetails.PreferencesLink,
		PrimaryColor:    SenderDetails.PrimaryColor,
	}

	outputHTML := bytes.NewBuffer([]byte{})

	err = htmlEmailTemplate.Execute(outputHTML, templateParams)
	if err != nil {
		panic(err)
	}

	outputText := bytes.NewBuffer([]byte{})

	err = textEmailTemplate.Execute(outputText, templateParams)
	if err != nil {
		panic(err)
	}

	return outputHTML.String(), outputText.String(), nil
}
