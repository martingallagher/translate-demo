package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/martingallagher/translate-demo/translate"
	"github.com/pkg/errors"
)

var httpClient = &http.Client{Timeout: time.Second * 3}

func translateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	var req *translate.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err == io.EOF {
		w.WriteHeader(http.StatusBadRequest)

		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error decoding JSON: %s", err)

		return
	}

	t := translate.NewTranslator(apiKey, httpClient)
	resp, err := t.Translate(req.Target, req.Text)

	if err != nil {
		if errors.Cause(err) == translate.ErrBadStatus {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		log.Printf("Translate error: %s", err)

		return
	}

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
