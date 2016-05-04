/*
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

package services

import (
	"time"

	"github.com/raiqub/jose/jws"
	"github.com/raiqub/jose/jws/data"
	"github.com/raiqub/jose/jwt"
)

// A Signer represents a service which provides token creation and signing.
type Signer struct {
	dtAdapter data.Signer
	config    Config
	keyCache  Cache
}

// NewSigner creates a new instance of Signer service.
func NewSigner(config Config, dt data.Signer) (*Signer, error) {
	dbKey, err := dt.GetKey(config.SignKeyID)
	if err != nil {
		return nil, err
	}
	rawKey, err := dbKey.Key()
	if err != nil {
		return nil, err
	}

	return &Signer{
		dt,
		config,
		Cache{*dbKey, rawKey},
	}, nil
}

// Create a new token and sign it.
func (s *Signer) Create(payload jwt.Claims) (string, error) {
	now := time.Now()

	payload.SetIssuer(s.config.Issuer)
	payload.SetExpireAt(now.Add(s.config.Duration))
	payload.SetNotBefore(now)
	payload.SetIssuedAt(now)

	header := &jws.RegisteredHeader{
		ID:        s.config.SignKeyID,
		Type:      jws.JWTHeaderType,
		Algorithm: s.keyCache.JWK.Algorithm,
		JWKSetURL: s.config.SetURL,
	}

	token := jws.SignedToken{
		Header:  header,
		Payload: payload,
	}
	out, err := token.EncodeAndSign(s.keyCache.RawKey)
	if err != nil {
		return "", err
	}

	return out, nil
}
