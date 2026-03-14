// Package models
package models

type File struct {
	Content string `json:"content"`
}

type ExecuteRequest struct {
	Language string `json:"language"`
	Version  string `json:"version"`
	Files    []File `json:"files"`
	Stdin    string `json:"stdin"`
}

type ExecuteResponse struct {
	Run struct {
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
		Code   int    `json:"code"`
	} `json:"run"`
}
