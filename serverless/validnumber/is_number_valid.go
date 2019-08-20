package validnumber

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/nexmo-community/nexmo-go"
)

type validResponse struct {
	IsValid bool   `json:"is_valid"`
	Error   string `json:"error,omitempty"`
}

type InsightRetriever interface {
	getInsight(phoneNumber string) (bool, error)
}

type NexmoInsightRetriever struct{}

func (nexmoRetriever NexmoInsightRetriever) getInsight(phoneNumber string) (bool, error) {
	auth := nexmo.NewAuthSet()
	auth.SetAPISecret(os.Getenv("NEXMO_API_KEY"), os.Getenv("NEXMO_API_SECRET"))
	client := nexmo.New(http.DefaultClient, auth)

	insight, _, err := client.Insight.GetAdvancedInsight(nexmo.AdvancedInsightRequest{
		Number: phoneNumber,
	})

	return insight.ValidNumber == nexmo.ValidNumberStatusValid, err
}

func IsNumberValid(w http.ResponseWriter, r *http.Request) {
	isNumberValid(w, r, NexmoInsightRetriever{})
}

func isNumberValid(w http.ResponseWriter, r *http.Request, ir InsightRetriever) {
	var request struct {
		PhoneNumber string `json:"phone_number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeResponse(w, validResponse{
			IsValid: false,
			Error:   err.Error(),
		})
		return
	}

	if request.PhoneNumber == "" {
		writeResponse(w, validResponse{
			IsValid: false,
			Error:   "You must provide a phone_number.",
		})
		return
	}

	if isValid, err := ir.getInsight(request.PhoneNumber); err != nil {
		writeResponse(w, validResponse{
			IsValid: isValid,
			Error:   err.Error(),
		})
	} else {
		writeResponse(w, validResponse{
			IsValid: isValid,
		})
	}
}

func writeResponse(writer http.ResponseWriter, response interface{}) {
	json.NewEncoder(writer).Encode(response)
}
