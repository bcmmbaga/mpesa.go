package mpesa

import (
	"errors"
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

	return &Application{
		Type:       apiType,
		Key:        applicationKey,
		SessionKey: "",
		market:     applicationMarket,
	}, nil
}
