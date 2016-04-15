package services

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/raiqub/jose/jwk"
	"github.com/raiqub/jose/jwk/keygen"
	"github.com/raiqub/jose/jwk/services"
	"github.com/raiqub/jose/jwt"
	"gopkg.in/mgo.v2"
	"gopkg.in/raiqub/eval.v0"
	"gopkg.in/raiqub/web.v0"
)

const (
	audience   = "testing"
	issuer     = "auth.example.com"
	duration   = time.Second * 10
	keyColName = "keys"
	keySize    = 2048
)

/*func benchmarkTokenCreation(kid string, b *testing.B) {
	config, err := createJWTConfig()
	if err != nil {
		b.Fatal(err)
	}
	config.SignKeyId = kid

	provider := NewSigner(*config)
	token, err := provider.Create(createJWTContext())
	if err != nil {
		b.Fatalf("Error creating token: %v", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := provider.Validate(audience, token)
		if err != nil {
			b.Fatalf("Error validating token: %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkECTokenCreation(b *testing.B) {
	benchmarkTokenCreation(eckid, b)
}

func BenchmarkRSATokenCreation(b *testing.B) {
	benchmarkTokenCreation(rsakid, b)
}

func benchmarkTokenValidation(kid string, b *testing.B) {
	config, err := createJWTConfig()
	if err != nil {
		b.Fatal(err)
	}
	config.SignKeyId = kid

	provider := NewSigner(*config)
	context := createJWTContext()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := provider.Create(context)
		if err != nil {
			b.Fatalf("Error creating token: %v", err)
		}
	}

	b.StopTimer()
}

func BenchmarkTokenECValidation(b *testing.B) {
	benchmarkTokenValidation(eckid, b)
}

func BenchmarkTokenRSAValidation(b *testing.B) {
	benchmarkTokenValidation(rsakid, b)
}*/

func testCreateAndValidate(
	alg string,
	session *mgo.Session,
	t *testing.T,
) {
	key, err := keygen.New(alg, keySize, 1)
	if err != nil {
		t.Fatalf("Error generating new key: %v", err)
	}

	if err := session.DB("").
		C(keyColName).
		Insert(key); err != nil {
		t.Fatalf("Error saving key to database: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			jwkset := jwk.Set{Keys: []jwk.Key{*key}}
			web.JSONWrite(w, http.StatusOK, jwkset)
		}))
	defer ts.Close()

	signer, err := NewSigner(Config{
		Issuer:    issuer,
		SetURL:    ts.URL,
		SignKeyId: key.ID,
		Duration:  duration,
	}, session.DB("").C(keyColName))
	if err != nil {
		t.Fatalf("Error creating signer: %v", err)
	}

	token, err := signer.Create(createJWTPayload())
	if err != nil {
		t.Errorf("Error creating token: %v", err)
	}

	if len(strings.TrimSpace(token)) == 0 {
		t.Error("Empty token")
	}

	cliJWKSet := services.NewSetClient(ts.URL)
	verifier, err := NewVerifier(cliJWKSet, issuer)
	if err != nil {
		t.Fatalf("Error creating verifier: %v", err)
	}
	if len(verifier.keys) == 0 {
		t.Fatal("No keys loaded from URL")
	}

	vToken, err := verifier.Verify(token)
	if err != nil {
		t.Fatalf("The token cannot be validated: %v", err)
	}
	if vToken == nil {
		t.Error("Invalid token")
	}

	//t.Logf("%s token: %s", alg, token)
}

func TestCreateAndValidate(t *testing.T) {
	env := eval.PrepareMongoDBEnvironment(t)
	if env == nil {
		return
	}
	defer env.Dispose()
	session := env.Session()

	// ECDSA
	testCreateAndValidate("ES256", session, t)
	testCreateAndValidate("ES384", session, t)
	testCreateAndValidate("ES512", session, t)

	// RSA
	testCreateAndValidate("RS256", session, t)
	testCreateAndValidate("RS384", session, t)
	testCreateAndValidate("RS512", session, t)
}

func createJWTPayload() *jwt.Payload {
	return &jwt.Payload{
		Audience: audience,
		Subject:  "gG26se5wyWDOEjaNHwlXm2i9G3mnYGbG62BBq3ZE",
		Scopes:   []string{"owner", "vehicle", "freight"},
		User: &jwt.PayloadUser{
			Name:    "John Doe",
			Email:   "john.doe@example.com",
			Country: "USA",
		},
	}
}