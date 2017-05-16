package translate

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pkg/errors"
)

type mockAPIClient struct {
	statusCode int
	err        error
	buf        []byte
}

func (c *mockAPIClient) Do(req *http.Request) (*http.Response, error) {
	if c.err != nil {
		return nil, c.err
	}

	return &http.Response{
		StatusCode: c.statusCode,
		Body:       ioutil.NopCloser(bytes.NewBuffer(c.buf)),
	}, nil
}

func newMockAPIClient(statusCode int, err error, buf []byte) APIClient {
	return &mockAPIClient{statusCode, err, buf}
}

func mockResponse(lang, text string) []byte {
	if lang == "" || text == "" {
		return nil
	}

	var (
		buf  = &bytes.Buffer{}
		data response
	)

	data.Data.Translations = append(data.Data.Translations, &responseTranslation{
		TranslatedText:         text,
		DetectedSourceLanguage: lang,
	})

	json.NewEncoder(buf).Encode(data)

	return buf.Bytes()
}

func TestTranslator(t *testing.T) {
	const apiKey = "N/A"

	errResp := errors.New("mock error")
	tests := []struct {
		statusCode                         int
		srcLang, dstLang, srcText, dstText string
		err, expectedErr                   error
	}{
		{200, "en", "de", "", "", nil, io.EOF},
		{200, "en", "de", "hello world", "Hallo Welt", nil, nil},
		{200, "zh-CN", "en", "你好，世界", "Hello world", nil, nil},
		{400, "00", "", "", "", nil, ErrBadStatus},
		{200, "en", "", "", "", errResp, errResp},
	}

	for _, v := range tests {
		c := newMockAPIClient(v.statusCode, v.err, mockResponse(v.srcLang, v.dstText))
		translator := NewTranslator(apiKey, c)
		s, err := translator.Translate(v.dstLang, v.srcText)

		if cause := errors.Cause(err); cause != v.expectedErr {
			t.Fatalf("Got %q, expected %q", cause, v.err)
		}

		// Skip further tests for failed calls
		if err != nil {
			continue
		}

		if s.DetectedLanguage != v.srcLang {
			t.Fatalf("Got %q, expected %q", s, v.dstText)
		} else if s.Translation != v.dstText {
			t.Fatalf("Got %q, expected %q", s.Translation, v.dstText)
		}
	}
}
