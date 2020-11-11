package mpesa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// APIEnviroment ...
type APIEnviroment string

// Market describe vodacom/vodafone operation territory
type Market string

const (
	baseURL = "https://openapi.m-pesa.com"

	// Sandbox dedicated enviroment for testing purposes.
	Sandbox APIEnviroment = "sandbox"

	// Production ready enviroment
	Production APIEnviroment = "openapi"

	// VodafoneGHANA country code GHA and currency code GHS
	VodafoneGHANA Market = "vodafoneGHA"

	// VodacomTanzania country code TZN and currency code TZS
	VodacomTanzania Market = "vodacomTZN"
)

// Application ...
type Application struct {
	client *http.Client

	// API enviroment for your application
	Type APIEnviroment

	// Unique Key needed to authorise and authenticate your application on the server.
	// created with the creation of a new application.
	Key string

	// Session Key acts as an access token that authorises the rest of your REST API calls to the system
	SessionKey string

	market Market
}

// NewApplication creates and returns new mpesa application
// you can pass an empty applicationKey as long as MPESA_APLICATION_KEY env has been set in your enviroment.
func NewApplication(applicationKey string, applicationMarket Market, apiType APIEnviroment) (*Application, error) {
	if apiType == "" || applicationMarket == "" {
		return nil, errors.New("Failed to create new application")
	}

	if applicationKey == "" {
		var key string

		if key = os.Getenv("MPESA_APPLICATION_KEY"); key == "" {
			return nil, errors.New("failed to create new application, application key is missing")
		}
		applicationKey = key
	}

	app := &Application{
		client:     &http.Client{},
		Type:       apiType,
		Key:        applicationKey,
		SessionKey: "",
		market:     applicationMarket,
	}

	if _, err := app.getSessionKey(); err != nil {
		return nil, err
	}

	return app, nil
}

// newRequest create new *http.Request with additional headers parameters required by MPESA API
func (app *Application) newRequest(method, url string, payload interface{}) (*http.Request, error) {

	var buf io.Reader

	if payload != nil {
		b, err := json.Marshal(&payload)
		if err != nil {
			return nil, err
		}

		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", app.SessionKey))
	req.Header.Set("Origin", "*")

	return req, nil
}

// send makes a request to the API, the response body will be unmarshaled into v
func (app *Application) send(req *http.Request, v interface{}) error {

	resp, err := app.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		data, err := ioutil.ReadAll(resp.Body)

		if err == nil && len(data) > 0 {
			json.Unmarshal(data, &v)
		}

		return errors.New("Handle this error")
	}

	if v != nil {

		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return err
		}

	}

	return nil
}
