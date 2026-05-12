package email

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
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
	Link:            "https://useabstractions.com",
	LogoDark:        "https://arterylabs-dev.s3.us-east-1.amazonaws.com/abstractions/logo-dark.png",
	LogoLight:       "https://arterylabs-dev.s3.us-east-1.amazonaws.com/abstractions/logo-light.png",
	LegalName:       "Artery Software Labs",
	Address:         "Bangalore, India",
	UnsubscribeLink: "https://useabstractions.com/unsubscribe",
	PreferencesLink: "https://useabstractions.com/preferences",
	PrimaryColor:    "#14b7d3",
	Email:           "support@useabstractions.com",
}

type SenderParams struct {
	Name            string
	Link            string
	LogoDark        string
	LogoLight       string
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
	SenderLogoDark  string
	SenderLogoLight string
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
		return "", "", fmt.Errorf("build email html template: %w", err)
	}

	textEmailTemplate, err := template.New("email").Parse(textTemplate)
	if err != nil {
		return "", "", fmt.Errorf("build email text template: %w", err)
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
		SenderLogo:      SenderDetails.LogoDark,
		SenderLogoDark:  SenderDetails.LogoDark,
		SenderLogoLight: SenderDetails.LogoLight,
		SenderLegalName: SenderDetails.LegalName,
		SenderAddress:   SenderDetails.Address,
		UnsubscribeLink: SenderDetails.UnsubscribeLink,
		PreferencesLink: SenderDetails.PreferencesLink,
		PrimaryColor:    SenderDetails.PrimaryColor,
	}

	outputHTML := bytes.NewBuffer([]byte{})

	err = htmlEmailTemplate.Execute(outputHTML, templateParams)
	if err != nil {
		return "", "", fmt.Errorf("execute html email template: %w", err)
	}

	outputText := bytes.NewBuffer([]byte{})

	err = textEmailTemplate.Execute(outputText, templateParams)
	if err != nil {
		return "", "", fmt.Errorf("execute text email template: %w", err)
	}

	return outputHTML.String(), outputText.String(), nil
}
