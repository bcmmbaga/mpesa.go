package mpesa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	publicKey = map[APIEnviroment]string{
		Sandbox:    `MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEArv9yxA69XQKBo24BaF/D+fvlqmGdYjqLQ5WtNBb5tquqGvAvG3WMFETVUSow/LizQalxj2ElMVrUmzu5mGGkxK08bWEXF7a1DEvtVJs6nppIlFJc2SnrU14AOrIrB28ogm58JjAl5BOQawOXD5dfSk7MaAA82pVHoIqEu0FxA8BOKU+RGTihRU+ptw1j4bsAJYiPbSX6i71gfPvwHPYamM0bfI4CmlsUUR3KvCG24rB6FNPcRBhM3jDuv8ae2kC33w9hEq8qNB55uw51vK7hyXoAa+U7IqP1y6nBdlN25gkxEA8yrsl1678cspeXr+3ciRyqoRgj9RD/ONbJhhxFvt1cLBh+qwK2eqISfBb06eRnNeC71oBokDm3zyCnkOtMDGl7IvnMfZfEPFCfg5QgJVk1msPpRvQxmEsrX9MQRyFVzgy2CWNIb7c+jPapyrNwoUbANlN8adU1m6yOuoX7F49x+OjiG2se0EJ6nafeKUXw/+hiJZvELUYgzKUtMAZVTNZfT8jjb58j8GVtuS+6TM2AutbejaCV84ZK58E2CRJqhmjQibEUO6KPdD7oTlEkFy52Y1uOOBXgYpqMzufNPmfdqqqSM4dU70PO8ogyKGiLAIxCetMjjm6FCMEA3Kc8K0Ig7/XtFm9By6VxTJK1Mg36TlHaZKP6VzVLXMtesJECAwEAAQ==`,
		Production: `MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAietPTdEyyoV/wvxRjS5pSn3ZBQH9hnVtQC9SFLgM9IkomEX9Vu9fBg2MzWSSqkQlaYIGFGH3d69Q5NOWkRo+Y8p5a61sc9hZ+ItAiEL9KIbZzhnMwi12jUYCTff0bVTsTGSNUePQ2V42sToOIKCeBpUtwWKhhW3CSpK7S1iJhS9H22/BT/pk21Jd8btwMLUHfVD95iXbHNM8u6vFaYuHczx966T7gpa9RGGXRtiOr3ScJq1515tzOSOsHTPHLTun59nxxJiEjKoI4Lb9h6IlauvcGAQHp5q6/2XmxuqZdGzh39uLac8tMSmY3vC3fiHYC3iMyTb7eXqATIhDUOf9mOSbgZMS19iiVZvz8igDl950IMcelJwcj0qCLoufLE5y8ud5WIw47OCVkD7tcAEPmVWlCQ744SIM5afw+Jg50T1SEtu3q3GiL0UQ6KTLDyDEt5BL9HWXAIXsjFdPDpX1jtxZavVQV+Jd7FXhuPQuDbh12liTROREdzatYWRnrhzeOJ5Se9xeXLvYSj8DmAI4iFf2cVtWCzj/02uK4+iIGXlX7lHP1W+tycLS7Pe2RdtC2+oz5RSSqb5jI4+3iEY/vZjSMBVk69pCDzZy4ZE8LBgyEvSabJ/cddwWmShcRS+21XvGQ1uXYLv0FCTEHHobCfmn2y8bJBb/Hct53BaojWUCAwEAAQ==`,
	}
)

type getSessionResp struct {
	// The response code for the transaction.
	Code string `json:"output_ResponseCode"`

	// The response description for the transaction.
	Description string `json:"output_ResponseDesc"`

	// The SessionKey that can be used to call other APIs.
	SessionID string `json:"output_SessionID"`
}

// getSession retrieve Session Key which authorises the rest of API calls to the system.
// Endpoint /[api_enviroment]/ipg/v2/[market]/getSession/
func (app *Application) getSessionKey() (string, error) {

	var sessionResp getSessionResp

	encryptedKey, err := encryptAPIKey(publicKey[app.Type], app.Key)
	if err != nil {
		return "", err
	}

	sessionEndpoint := fmt.Sprintf("%s/%s/ipg/v2/%s/getSession/", baseURL, app.Type, app.market)

	req, err := http.NewRequest(http.MethodGet, sessionEndpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", encryptedKey))
	req.Header.Set("Origin", "*")

	resp, err := app.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&sessionResp); err != nil {
		return "", err
	}

	app.SessionKey = sessionResp.SessionID

	return sessionResp.SessionID, nil
}

func encryptAPIKey(pbKey string, APIkey string) (string, error) {

	base64Str, err := base64.StdEncoding.DecodeString(pbKey)
	if err != nil {
		// just for now no better error handling
		return "", errors.New("Handle this error")
	}

	parsedPublicKey, err := x509.ParsePKIXPublicKey(base64Str)
	if err != nil {
		return "", errors.New("Handle this error")
	}

	pk, ok := parsedPublicKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("Handle this error")
	}

	encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, pk, []byte(APIkey))
	if err != nil {
		return "", errors.New("Handle this error")
	}

	return base64.StdEncoding.EncodeToString(encryptedKey), nil
}
