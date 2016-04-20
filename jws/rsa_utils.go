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

package jws

import (
	"crypto/x509"
	"encoding/pem"
)

var (
	rsaParsers = []func([]byte) (interface{}, error){
		func(key []byte) (interface{}, error) {
			return x509.ParsePKCS1PrivateKey(key)
		},
		func(key []byte) (interface{}, error) {
			return x509.ParsePKCS8PrivateKey(key)
		},
		func(key []byte) (interface{}, error) {
			return x509.ParsePKIXPublicKey(key)
		},
		func(key []byte) (interface{}, error) {
			cert, err := x509.ParseCertificate(key)
			if err != nil {
				return nil, err
			}

			return cert.PublicKey, nil
		},
	}
)

// ParseRSAFromPEM decodes PEM encoded PKCS1 or PKCS8.
func ParseRSAFromPEM(key []byte) (interface{}, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded(0)
	}

	var parsedKey interface{}
	for _, v := range rsaParsers {
		if parsedKey, err = v(block.Bytes); err == nil {
			return parsedKey, nil
		}
	}

	return nil, ErrParsingFromPEM(0)
}
