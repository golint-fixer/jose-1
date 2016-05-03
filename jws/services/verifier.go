package services

import (
	"errors"

	jwkservices "github.com/raiqub/jose/jwk/services"
	"github.com/raiqub/jose/jws"
	"github.com/raiqub/jose/jwt"
	"gopkg.in/raiqub/dot.v1"
)

// A Verifier represents a service which provides token decoding and validation.
type Verifier struct {
	issuer string
	keys   map[string]*Cache
}

// NewVerifier creates a new instance of Verifier service.
func NewVerifier(
	svcJWKSet jwkservices.SetService,
	issuer string,
) (*Verifier, error) {
	jwkset, jerr := svcJWKSet.GetCerts()
	if jerr != nil {
		// TODO proper error type
		return nil, errors.New(jerr.Type + ": " + jerr.Message)
	}

	result := &Verifier{
		issuer,
		make(map[string]*Cache, 0),
	}

	for _, k := range jwkset.Keys {
		rawKey, err := k.Key()
		if err != nil {
			return nil, err
		}

		result.keys[k.ID] = &Cache{k, rawKey}
		// TODO log loaded key
		//fmt.Println("[Verifier] Loaded key:", k.ID)
	}
	// TODO log loaded keys
	//fmt.Printf("[Verifier] Loaded %d keys\n", len(result.keys))

	return result, nil
}

// Verify specified token and decode it.
func (v *Verifier) Verify(rawtoken string) (*jws.SignedToken, error) {
	token, err := jws.DecodeAndValidate(
		rawtoken, nil, nil,
		func(header jws.Header) (interface{}, error) {
			var key *Cache
			var ok bool

			if key, ok = v.keys[header.GetID()]; !ok {
				return nil, ErrInvalidKeyID(header.GetID())
			}
			if header.GetAlgorithm() != key.JWK.Algorithm {
				return nil, ErrUnexpectedAlg(header.GetAlgorithm())
			}

			return key.RawKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Validate() ||
		token.Payload.GetIssuer() != v.issuer {
		return nil, ErrInvalidToken(0)
	}

	return token, nil
}

// VerifyScopes validates client and user scopes when available.
func (v *Verifier) VerifyScopes(
	claims jwt.ClientUserScopes,
	client, user []string,
) bool {
	var scopes []string

	if client != nil && len(client) > 0 {
		scopes = claims.GetScopes()
		if scopes == nil || len(scopes) == 0 {
			return false
		}

		if !dot.StringSlice(scopes).
			ExistsAny(client, false) {
			return false
		}
	}

	if user != nil && len(user) > 0 {
		scopes = claims.GetUserScopes()
		if scopes == nil || len(scopes) == 0 {
			return false
		}

		if !dot.StringSlice(scopes).
			ExistsAny(user, false) {
			return false
		}
	}

	return true
}