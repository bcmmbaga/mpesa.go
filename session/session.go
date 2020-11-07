/*
 * Copyright 2020 Infolabs Inc & Associates
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package session

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/mobilemoney/mpesa/pkg/errors"
)

var (

	ErrDecodePubKey = errors.New("error occurred while decoding public key")
)

type Application struct {

	//Name of the application is a human-readable name of
	//the application
	Name string `json:"name"`

	//Version number of the application allowing changes in API products
	//to be managed in different versions
	Version string `json:"version"`

	//Desc Free text to describe the use of application
	Desc string `json:"desc"`

	//APIKey - a unique authorization key used to authenticate the application
	//on the first call.
	//API Keys need to be encrypted in the 1st "Generate Session API call" to
	//create a valid session key to be used as an Access token for future calls.
	//Encrypting the APIKey is done by EncryptAPIKey
	APIKey string `json:"api_key"`

	//SessionLifeTime - Session Key has a finite lifetime of availability that can
	//be configured. Once it has expired, session is no longer usable and the caller
	//will need to authenticate again.
	SessionLifeTime int `json:"session_life_time"`

	//TrustedSources The originating caller can be limited to specific IP Addresses
	//as an additional security measure.
	TrustedSources []string `json:"trusted_sources"`
	
	
	PublicKey string `json:"public_key"`


	// todo: Scope
}

var _ Session = (*Application)(nil)

type Config struct {
	Application
}

type Session interface {
	GenerateSessionKey(ctx context.Context)

	//EncryptAPIKey
	//Log in to the OPENAPI portal with dev account. Create New Application
	//A new unique APIKey will be generated for the newly created application
	//copy and save the api key in the configuration file along side other
	//Application attributes like Application.Name and Application.Version
	//Copy the public Key from the web.
	//Steps in Encrypting the APIKey
	//1. Generate a decoded Base64 string from the Public Key
	//2. Generate an instance of an RSA cipher and use Base64 as the input
	//3. Encode the APIKey with RSA cipher and digest as Base64 string format
	// Now step (3) provides encrypted api key
	EncryptAPIKey()error
}


// New instantiates the users service implementation
func New(cfg Config) Session{

	name := cfg.Name
	version := cfg.Version
	desc := cfg.Desc
	key := cfg.APIKey

	return &Application{
		Name:            name,
		Version:         version,
		Desc:            desc,
		APIKey:          key,
		SessionLifeTime: 0,
		TrustedSources:  nil,
	}
}


func (a Application) GenerateSessionKey(ctx context.Context) {
	panic("implement me")
}

func (a Application) EncryptAPIKey()error {


	//Decode Pub Key
	key,err := base64.StdEncoding.DecodeString(a.PublicKey)

	if err != nil {
		return errors.Wrap(ErrDecodePubKey,err)
	}

	//get string from bytes
	keyStr := string(key)


	se := base64.StdEncoding.EncodeToString([]byte(keyStr))
	fmt.Println(se)                                        // YSBzdHJpbmc=

	sd, e := base64.StdEncoding.DecodeString(se)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(string(sd))
	panic("implement me")
}