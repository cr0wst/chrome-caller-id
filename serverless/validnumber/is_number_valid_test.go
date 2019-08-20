package validnumber

import (
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockNexmoValid struct{}
type mockNexmoInValid struct{}
type mockNexmoInValidWithError struct{}

func (nexmo mockNexmoValid) getInsight(phoneNumber string) (bool, error) {
	return true, nil
}

func (nexmo mockNexmoInValid) getInsight(phoneNumber string) (bool, error) {
	return false, nil
}

func (nexmo mockNexmoInValidWithError) getInsight(phoneNumber string) (bool, error) {
	return false, errors.New("An Error has Occurred")
}

func TestIsNumberValid(t *testing.T) {
	tests := []struct {
		body string
		ir   InsightRetriever
		want string
	}{
		{body: `{}`, ir: mockNexmoInValid{}, want: "{\"is_valid\":false,\"error\":\"You must provide a phone_number.\"}\n"},
		{body: ``, ir: mockNexmoInValid{}, want: "{\"is_valid\":false,\"error\":\"EOF\"}\n"},
		{body: `{"phone_number": "_"}`, ir: mockNexmoInValidWithError{}, want: "{\"is_valid\":false,\"error\":\"An Error has Occurred\"}\n"},
		{body: `{"phone_number": "1"}`, ir: mockNexmoInValid{}, want: "{\"is_valid\":false}\n"},
		{body: `{"phone_number": "15555555555"}`, ir: mockNexmoInValid{}, want: "{\"is_valid\":false}\n"},
		{body: `{"phone_number": "18885537704"}`, ir: mockNexmoValid{}, want: "{\"is_valid\":true}\n"},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/", strings.NewReader(test.body))
		req.Header.Add("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		isNumberValid(rec, req, test.ir)

		out, err := ioutil.ReadAll(rec.Result().Body)

		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		if got := string(out); got != test.want {
			t.Errorf("IsNumberValid(%q) = %q, want %q", test.body, got, test.want)
		}
	}
}
