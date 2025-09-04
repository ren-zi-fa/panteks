package validator

import (
	"bytes"
	"errors"

	"github.com/PuerkitoBio/goquery"
)

func ValidateSelector(body []byte, selector string) error {
    if selector == "" {
        return errors.New("selector cannot be empty")
    }

    doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
    if err != nil {
        return err
    }

    if doc.Find(selector).Length() == 0 {
        return errors.New("selector not found")
    }

    return nil
}