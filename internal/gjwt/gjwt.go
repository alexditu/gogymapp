package gjwt

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type GoogleOIDClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	jwt.RegisteredClaims
}

// taken from: https://www.googleapis.com/oauth2/v3/certs
var googleKeys = map[string]string{
	"6f7254101f56e41cf35c9926de84a2d552b4c6f1": "oUriU8GqbRw-avcMn95DGW1cpZR1IoM6L7krfrWvLSSCcSX6Ig117o25Yk7QWBiJpaPV0FbP7Y5-DmThZ3SaF0AXW-3BsKPEXfFfeKVc6vBqk3t5mKlNEowjdvNTSzoOXO5UIHwsXaxiJlbMRalaFEUm-2CKgmXl1ss_yGh1OHkfnBiGsfQUndKoHiZuDzBMGw8Sf67am_Ok-4FShK0NuR3-q33aB_3Z7obC71dejSLWFOEcKUVCaw6DGVuLog3x506h1QQ1r0FXKOQxnmqrRgpoHqGSouuG35oZve1vgCU4vLZ6EAgBAbC0KL35I7_0wUDSMpiAvf7iZxzJVbspkQ",
	"b9ac601d131fd4ffd556ff032aab188880cde3b9": "trD9XzkQVbaVs5NeV-PrHMYGm9JsfXKKoPJWuU8zcA5T7sp25j4KvJAPgSdFO1x6AiVtxwKGUBnsr9gnNaiSM3qs_1_iJT09E_iqHUyaTiqf4wkEHA5ABinBkORsjQzZajsbbhtkv4Yw4vF44g2WhchdjLThpBB96px-RV4C0ZK8beA-4cNEYhybYBsEjYDZLWAIKxtt-ZNc01AhM1p5nIDLp6Z05hAJBVazj7Ac3JT_CwgYlY3MvLZSJIQjOZwBRLNl9wJhewiNvfIH3ijbPVKzLEyt5toNqsSyuBZtLr-z4UKv2gsoKFSU-KdkRBnO3ZtqVYIsiZ-09IEN1pL33Q",
}

// taken from: https://www.googleapis.com/oauth2/v1/certs
var googleCertsPem = map[string]string{
	"b9ac601d131fd4ffd556ff032aab188880cde3b9": "-----BEGIN CERTIFICATE-----\nMIIDJzCCAg+gAwIBAgIJAND4bIHLtP7BMA0GCSqGSIb3DQEBBQUAMDYxNDAyBgNV\nBAMMK2ZlZGVyYXRlZC1zaWdub24uc3lzdGVtLmdzZXJ2aWNlYWNjb3VudC5jb20w\nHhcNMjMwOTI1MDQzNzU5WhcNMjMxMDExMTY1MjU5WjA2MTQwMgYDVQQDDCtmZWRl\ncmF0ZWQtc2lnbm9uLnN5c3RlbS5nc2VydmljZWFjY291bnQuY29tMIIBIjANBgkq\nhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtrD9XzkQVbaVs5NeV+PrHMYGm9JsfXKK\noPJWuU8zcA5T7sp25j4KvJAPgSdFO1x6AiVtxwKGUBnsr9gnNaiSM3qs/1/iJT09\nE/iqHUyaTiqf4wkEHA5ABinBkORsjQzZajsbbhtkv4Yw4vF44g2WhchdjLThpBB9\n6px+RV4C0ZK8beA+4cNEYhybYBsEjYDZLWAIKxtt+ZNc01AhM1p5nIDLp6Z05hAJ\nBVazj7Ac3JT/CwgYlY3MvLZSJIQjOZwBRLNl9wJhewiNvfIH3ijbPVKzLEyt5toN\nqsSyuBZtLr+z4UKv2gsoKFSU+KdkRBnO3ZtqVYIsiZ+09IEN1pL33QIDAQABozgw\nNjAMBgNVHRMBAf8EAjAAMA4GA1UdDwEB/wQEAwIHgDAWBgNVHSUBAf8EDDAKBggr\nBgEFBQcDAjANBgkqhkiG9w0BAQUFAAOCAQEAdL2D4ZVRxBt2TohXV+JpDFFZ92xH\nQH0OJ0bhbrfCc6AGBXx13IiLUwHok4jNZ0x+ZXQyDR9rKOdo5iTn4kQKD2blor5m\nj4r8aK/nXIxU7foxK0H7dJMALMdslAl5L3LKrE5beNLk/v2kfQM0pTqzeGaqTzNg\n3ZSHPQgJ/i8ES8+7dV12A+ct4nv3DT1M6rmCa6+AowlzpIllBlOkIe45qamPfY4j\niyNlfOlKA+r7JE/vXQeMHTucNSBzHlEL48wbo4nS4ftnEoI/9E/ryXZ4zp9YVDnU\nXQ2jfVgymwzNjPSoFJa+BAJOsx3Ig3Hor322OSXgNF3sIAb1xu4598k08g==\n-----END CERTIFICATE-----\n",
	"6f7254101f56e41cf35c9926de84a2d552b4c6f1": "-----BEGIN CERTIFICATE-----\nMIIDJjCCAg6gAwIBAgIIYwnpFReDI4QwDQYJKoZIhvcNAQEFBQAwNjE0MDIGA1UE\nAwwrZmVkZXJhdGVkLXNpZ25vbi5zeXN0ZW0uZ3NlcnZpY2VhY2NvdW50LmNvbTAe\nFw0yMzA5MTcwNDM3NThaFw0yMzEwMDMxNjUyNThaMDYxNDAyBgNVBAMMK2ZlZGVy\nYXRlZC1zaWdub24uc3lzdGVtLmdzZXJ2aWNlYWNjb3VudC5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQChSuJTwaptHD5q9wyf3kMZbVyllHUigzov\nuSt+ta8tJIJxJfoiDXXujbliTtBYGImlo9XQVs/tjn4OZOFndJoXQBdb7cGwo8Rd\n8V94pVzq8GqTe3mYqU0SjCN281NLOg5c7lQgfCxdrGImVsxFqVoURSb7YIqCZeXW\nyz/IaHU4eR+cGIax9BSd0qgeJm4PMEwbDxJ/rtqb86T7gVKErQ25Hf6rfdoH/dnu\nhsLvV16NItYU4RwpRUJrDoMZW4uiDfHnTqHVBDWvQVco5DGeaqtGCmgeoZKi64bf\nmhm97W+AJTi8tnoQCAEBsLQovfkjv/TBQNIymIC9/uJnHMlVuymRAgMBAAGjODA2\nMAwGA1UdEwEB/wQCMAAwDgYDVR0PAQH/BAQDAgeAMBYGA1UdJQEB/wQMMAoGCCsG\nAQUFBwMCMA0GCSqGSIb3DQEBBQUAA4IBAQCaURaabcRD7dB7spKZgBjd17O3e0fK\nXa37RPvoEVYKBGyzuQCw1nWJhGKdQHnDXxSt9UTqOZLkVAfTC8/DAZdCnY6rJG01\nuP1gk77PxqBnfOMaBQNLO4wm1c/XkTtpPDvPJ3Pwy4xvd1pr2YMz2SDj4yNuZIaB\njiREcQoEbYlvcfD1e/mxrSOErogwZoZcNdo+QFKQPnODvJvTHXDtAoqT4/Q96ynC\n8hBx+F9PQvvlVUMl4tuQvtMjpkv58hG4gPovdzkadOOENvmObgjiTL6Klj+asAVq\ngTUd53NSpW/grtxob7s/x//B1jyNq5C2iebVM7n4pO39SOmiHVhAm1Eb\n-----END CERTIFICATE-----\n",
}

// TODO: fetch keys/certs from: https://www.googleapis.com/oauth2/v3/certs
// according to docs they are updated regulary (see https://developers.google.com/identity/gsi/web/guides/verify-google-id-token)
func UpdateGoogleKeys() {

}

func ParseToken(tokenStr string) (*jwt.Token, *GoogleOIDClaims, error) {
	var c GoogleOIDClaims
	token, err := jwt.ParseWithClaims(tokenStr, &c, func(token *jwt.Token) (interface{}, error) {

		kid, ok := token.Header["kid"]
		if !ok {
			return nil, fmt.Errorf("missing 'kid' header from token")
		}

		kidStr := kid.(string)

		//
		// TODO: read about JWK token and how to parse them
		//
		// keyEnc, ok := googleKeys[kidStr]
		// if !ok {
		// 	return nil, fmt.Errorf("missing key with kid: '%s'", kidStr)
		// }

		// key, err := b64.URLEncoding.DecodeString(keyEnc)
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to decode base64 url key with kid '%s': %s", kidStr, err)
		// }
		//

		pemCert, ok := googleCertsPem[kidStr]
		if !ok {
			return nil, fmt.Errorf("missing key with kid: '%s'", kidStr)
		}

		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pemCert))
		if err != nil {
			return nil, fmt.Errorf("failed to parse PEM key with kid '%s': %s", kidStr, err)
		}

		return key, nil
	})

	if err != nil {
		return nil, nil, err
	}

	return token, &c, nil
}

// aud is the client_id received from Google when registered the app
// based on https://developers.google.com/identity/gsi/web/guides/verify-google-id-token
func ValidateToken(t *jwt.Token, clientId string) error {
	if t == nil {
		return fmt.Errorf("nil token")
	}

	// 1. value of aud in the ID token is equal to one of your app's client ID
	aud, err := t.Claims.GetAudience()
	if err != nil {
		return fmt.Errorf("failed to GetAudience(): %s", err)
	}

	if len(aud) < 1 {
		return fmt.Errorf("token does not contain aud field")
	}

	if aud[0] != clientId {
		return fmt.Errorf("invalid audience: '%s'", aud)
	}

	// 2. value of iss in the ID token is equal to accounts.google.com or https://accounts.google.com
	iss, err := t.Claims.GetIssuer()
	if err != nil {
		return fmt.Errorf("failed to GetIssuer(): %s", err)
	}

	if iss != "accounts.google.com" && iss != "https://accounts.google.com" {
		return fmt.Errorf("invalid issuer: '%s'", iss)
	}

	// 3. expiry time (exp) of the ID token has not passed
	exp, err := t.Claims.GetExpirationTime()
	if err != nil || exp == nil {
		return fmt.Errorf("failed to GetExpirationTime(): %s", err)
	}

	if exp.Before(time.Now()) {
		return fmt.Errorf("token expired at %s", *exp)
	}

	return nil
}

func ValidateGoogleJWT(c *gin.Context, clientId string) (*GoogleOIDClaims, error) {
	err := verifyCSRF(c)
	if err != nil {
		return nil, err
	}

	credsToken := c.PostForm("credential")
	if credsToken == "" {
		return nil, fmt.Errorf("missing credential post form")
	}

	t, claims, err := ParseToken(credsToken)
	if err != nil {
		return nil, err
	}

	if t == nil {
		return nil, fmt.Errorf("token is nil")
	}

	if claims == nil {
		return nil, fmt.Errorf("claims is nil")
	}

	err = ValidateToken(t, clientId)
	if err != nil {
		// although invalid, the claims object is parsed successuflly
		return claims, err
	}

	return claims, nil

}

func verifyCSRF(c *gin.Context) error {
	csrfTokenCookie, err := c.Cookie("g_csrf_token")
	if err != nil {
		return err
	}

	csrfBodyCookie := c.PostForm("g_csrf_token")
	if csrfBodyCookie == "" {
		return fmt.Errorf("missing g_csrf_token post form")
	}

	if csrfTokenCookie != csrfBodyCookie {
		return fmt.Errorf("csrfTokenCookie != csrfBodyCookie (%s != %s)", csrfTokenCookie, csrfBodyCookie)
	}

	return nil
}
