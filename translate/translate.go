package translate

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const googleTranslateAPI = "https://translation.googleapis.com/language/translate/v2?key="

var (
	// ErrBadStatus - bad HTTP status code.
	ErrBadStatus = errors.New("bad HTTP status code")
	// ErrEmptyTranslations - empty translations.
	ErrEmptyTranslations = errors.New("empty translations")
)

// Translator defines the translator interface.
type Translator interface {
	Translate(lang, text string) (*Translation, error)
}

// APIClient defines a simple interface for the Do() method.
type APIClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type translator struct {
	endpoint string
	client   APIClient
}

// Translation represents a translation result.
type Translation struct {
	DetectedLanguage string `json:"detectedLanguage"`
	Translation      string `json:"translation"`
}

// Request represents a translation request.
type Request struct {
	Text   string `json:"q"`
	Target string `json:"target"`
}

type response struct {
	Data struct {
		Translations []*responseTranslation `json:"translations"`
	} `json:"data"`
}

type responseTranslation struct {
	DetectedSourceLanguage string `json:"detectedSourceLanguage"`
	TranslatedText         string `json:"translatedText"`
}

// NewTranslator returns a new translator with the given API key.
func NewTranslator(apiKey string, client APIClient) Translator {
	return &translator{
		endpoint: googleTranslateAPI + apiKey,
		client:   client,
	}
}

func (t *translator) Translate(lang, text string) (*Translation, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(Request{
		Text:   text,
		Target: lang,
	})

	if err != nil {
		// Memory exhausted
		return nil, err
	}

	r := ioutil.NopCloser(buf)
	req, err := http.NewRequest(http.MethodPost, t.endpoint, r)

	if err != nil {
		// Impossible unless invalid HTTP method provided
		return nil, err
	}

	resp, err := t.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(ErrBadStatus, "status code: %d", resp.StatusCode)
	}

	var translation *response

	if err = json.NewDecoder(resp.Body).Decode(&translation); err != nil {
		return nil, err
	}

	translations := translation.Data.Translations

	if len(translations) == 0 {
		return nil, ErrEmptyTranslations
	}

	return &Translation{
		DetectedLanguage: translations[0].DetectedSourceLanguage,
		Translation:      translations[0].TranslatedText,
	}, nil
}
