package models

type SendEmailRequest struct {
	Destination Destination `json:"Destination" binding:"required"`
	Message     Message     `json:"Message" binding:"required"`
}

type Destination struct {
	ToAddresses  []string `json:"ToAddresses" binding:"required"`
	CcAddresses  []string `json:"CcAddresses"`
	BccAddresses []string `json:"BccAddresses"`
}

type Html struct {
	Data    string `json:"Data"`
	Charset string `json:"Charset"`
}

type Body struct {
	Text Html `json:"Text"`
	Html Html `json:"Html"`
}

type Subject struct {
	Data    string `json:"Data" binding:"required"`
	Charset string `json:"Charset"`
}

type Message struct {
	Subject Subject `json:"Subject" binding:"required"`
	Body    Body    `json:"Body" binding:"required"`
}
type SendEmailResponse struct {
	MessageId string `json:"MessageId"`
	RequestId string `json:"RequestId"`
}

type SESError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}
