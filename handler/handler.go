package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rayou/go-theperfectgiftcard"
)

type requestPayload struct {
	CardNo string `json:"card_no"`
	Pin    string `json:"pin"`
}

// Handler is the main strcut for http server handlers.
type Handler struct {
	*http.ServeMux
	client httpClient
}

// The Perfect Gift Card HTTP Client Interface
type httpClient interface {
	GetCard(cardNo string, pin string) (*theperfectgiftcard.Card, *theperfectgiftcard.Response, error)
}

// NewHandler creates a handler for http server
func NewHandler(client httpClient) *Handler {
	mux := http.NewServeMux()
	h := &Handler{ServeMux: mux, client: client}
	h.HandleFunc("/card", h.cardHandler)
	return h
}

func (h *Handler) cardHandler(w http.ResponseWriter, r *http.Request) {
	var p requestPayload
	json.NewDecoder(r.Body).Decode(&p)

	if p.CardNo == "" {
		respondWithError(w, http.StatusBadRequest, "card no is required")
		return
	}

	if p.Pin == "" {
		respondWithError(w, http.StatusBadRequest, "pin is required")
		return
	}

	card, resp, err := h.client.GetCard(p.CardNo, p.Pin)
	if err != nil {
		errMsg := err.Error()
		if resp != nil && resp.StatusCode == http.StatusInternalServerError {
			errMsg = "the perfect gift card website unreachable."
		}
		respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}
	respondWithJSON(w, resp.StatusCode, &card)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
