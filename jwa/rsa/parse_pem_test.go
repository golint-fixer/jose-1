/*
 * Copyright 2012 Dave Grijalva
 * Copyright 2016 Fabrício Godoy
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rsa_test

import (
	"io/ioutil"
	"testing"

	"github.com/raiqub/jose/jwa/rsa"
)

const (
	privKeyFile = "test/sample_key"
	pubKeyFile  = "test/sample_key.pub"
)

func TestRSAKeyParsing(t *testing.T) {
	key, _ := ioutil.ReadFile(privKeyFile)
	pubKey, _ := ioutil.ReadFile(pubKeyFile)
	badKey := []byte("All your base are belong to key")

	if _, e := rsa.ParseFromPEM(key); e != nil {
		t.Errorf("Failed to parse valid private key: %v", e)
	}

	if _, e := rsa.ParseFromPEM(pubKey); e != nil {
		t.Errorf("Failed to parse valid public key: %v", e)
	}

	if k, e := rsa.ParseFromPEM(badKey); e == nil {
		t.Errorf("Parsed invalid key as valid private key: %v", k)
	}
}
