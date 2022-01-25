package model

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	TitleInvalidValue  = errors.New("title had empty value, expect a valid one")
	SigneeInvalidValue = errors.New("signee had empty value, expect a valid one")
)

type Document struct {
	ID      int64   `json:"id"`
	Title   string  `json:"title"`
	Content Content `json:"content"`
	Signee  string  `json:"signee"`
}

type Content struct {
	Header string `json:"header"`
	Data   string `json:"data"`
}

func (c *Content) Scan(src interface{}) error {
	val := src.([]uint8)
	return json.Unmarshal(val, &c)
}

func (d *Document) Validate() error {
	d.Title = strings.TrimSpace(d.Title)
	d.Signee = strings.TrimSpace(d.Signee)
	if d.Title == "" {
		return TitleInvalidValue
	}
	if d.Signee == "" {
		return SigneeInvalidValue
	}
	return nil
}
