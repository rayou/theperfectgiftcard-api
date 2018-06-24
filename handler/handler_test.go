package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gocolly/colly"

	tpgc "github.com/rayou/go-theperfectgiftcard"
)

type mockHTTPClient struct {
	t     *testing.T
	card  *tpgc.Card
	resp  *tpgc.Response
	error error
}

func (m *mockHTTPClient) GetCard(cardNo string, pin string) (*tpgc.Card, *tpgc.Response, error) {
	testValue(m.t, cardNo, expectedCard.CardNo)
	testValue(m.t, pin, "0000")
	return m.card, m.resp, m.error
}

var bodyTemplate = "{\"card_no\": \"%s\", \"pin\": \"%s\"}"
var expectedCard = &tpgc.Card{
	CardNo:           "50211234567890",
	AccountNo:        "000000000",
	LoadsToDate:      "$100.00",
	PurchasesToDate:  "-$54.32",
	AvailableBalance: "$12.34",
	PurchasedDate:    "1 Jan 2018",
	ExpiryDate:       "1 Jan 2021",
	Transactions: []tpgc.Transaction{
		{
			Date:        "1 Jan 2018 12:04:45 PM",
			Details:     "Store Address",
			Description: "Refund - Store Address",
			Amount:      "$100.00",
			Balance:     "$100.00",
		},
		{
			Date:        "2 Jan 2018 07:50:53 PM",
			Details:     "Store A",
			Description: "Purchase - Store A",
			Amount:      "$12.34-",
			Balance:     "$56.78",
		},
	},
}

func testValue(t *testing.T, got interface{}, want interface{}) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Expected: %v, got: %v", want, got)
	}
}

func buildMockHTTPClient(t *testing.T, card *tpgc.Card, statusCode int, err error) *mockHTTPClient {
	return &mockHTTPClient{
		t:    t,
		card: card,
		resp: &tpgc.Response{
			Response: &colly.Response{
				StatusCode: statusCode,
			},
		},
		error: err,
	}
}

func TestHandler(t *testing.T) {
	client := buildMockHTTPClient(t, expectedCard, http.StatusOK, nil)
	s := httptest.NewServer(NewHandler(client))

	resp, err := http.Post(s.URL+"/card", "application/json", strings.NewReader(fmt.Sprintf(bodyTemplate, expectedCard.CardNo, "0000")))
	testValue(t, err == nil, true)
	defer resp.Body.Close()

	var m tpgc.Card
	err = json.NewDecoder(resp.Body).Decode(&m)
	testValue(t, err == nil, true)
	testValue(t, m.CardNo, expectedCard.CardNo)
	testValue(t, m.AvailableBalance, expectedCard.AvailableBalance)
	testValue(t, m.Transactions[0].Balance, expectedCard.Transactions[0].Balance)
}

func TestHandlerMissingCardNo(t *testing.T) {
	client := buildMockHTTPClient(t, nil, http.StatusBadRequest, nil)
	s := httptest.NewServer(NewHandler(client))

	resp, err := http.Post(s.URL+"/card", "application/json", strings.NewReader(""))
	testValue(t, err == nil, true)
	defer resp.Body.Close()

	var m map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	testValue(t, err == nil, true)
	testValue(t, resp.StatusCode, http.StatusBadRequest)
	testValue(t, m["error"], "card no is required")
}

func TestHandlerMissingPin(t *testing.T) {
	client := buildMockHTTPClient(t, nil, http.StatusBadRequest, nil)
	s := httptest.NewServer(NewHandler(client))

	resp, err := http.Post(s.URL+"/card", "application/json",
		strings.NewReader(fmt.Sprintf(bodyTemplate, expectedCard.CardNo, "")))
	testValue(t, err == nil, true)
	defer resp.Body.Close()

	var m map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	testValue(t, err == nil, true)
	testValue(t, resp.StatusCode, http.StatusBadRequest)
	testValue(t, m["error"], "pin is required")
}

func TestClientError(t *testing.T) {
	errMsg := "this is a client error"
	client := &mockHTTPClient{
		t:     t,
		error: errors.New(errMsg),
	}
	s := httptest.NewServer(NewHandler(client))

	resp, err := http.Post(s.URL+"/card", "application/json",
		strings.NewReader(fmt.Sprintf(bodyTemplate, expectedCard.CardNo, "0000")))
	testValue(t, err == nil, true)
	defer resp.Body.Close()

	var m map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	testValue(t, err == nil, true)
	testValue(t, resp.StatusCode, http.StatusBadRequest)
	testValue(t, m["error"], errMsg)
}

func TestTargetPageBadRequest(t *testing.T) {
	errMsg := "Invalid card number or password"
	client := buildMockHTTPClient(t, nil, http.StatusBadRequest, errors.New(errMsg))
	s := httptest.NewServer(NewHandler(client))

	resp, err := http.Post(s.URL+"/card", "application/json",
		strings.NewReader(fmt.Sprintf(bodyTemplate, expectedCard.CardNo, "0000")))
	testValue(t, err == nil, true)
	defer resp.Body.Close()

	var m map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	testValue(t, err == nil, true)
	testValue(t, resp.StatusCode, http.StatusBadRequest)
	testValue(t, m["error"], errMsg)
}

func TestTargetPageInternalError(t *testing.T) {
	errMsg := "internal server error"
	client := buildMockHTTPClient(t, nil, http.StatusInternalServerError, errors.New(errMsg))
	s := httptest.NewServer(NewHandler(client))

	resp, err := http.Post(s.URL+"/card", "application/json",
		strings.NewReader(fmt.Sprintf(bodyTemplate, expectedCard.CardNo, "0000")))
	testValue(t, err == nil, true)
	defer resp.Body.Close()

	var m map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	testValue(t, err == nil, true)
	testValue(t, resp.StatusCode, http.StatusBadRequest)
	testValue(t, m["error"], "the perfect gift card website unreachable.")
}
