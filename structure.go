package main

import "encoding/xml"

type Config struct {
	XMLName         xml.Name `xml:"Config"`
	APIAddress      string   `xml:"APIAddress"`
	RedirectAddress string   `xml:"RedirectAddress"`
}

type ErrorCode struct {
	Name        string   `json:"name"`
	Regex       string   `json:"regex"`
	Card        string   `json:"card"`
	Comment     string   `json:"comment"`
	Description []string `json:"description"`
}

type ErrorResponse struct {
	Error    string      `json:"error"`
	Found    int         `json:"found"`
	InfoList []ErrorInfo `json:"infolist"`
}

type ErrorInfo struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Info string `json:"info"`
}
