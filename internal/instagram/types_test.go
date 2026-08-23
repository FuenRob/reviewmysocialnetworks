package instagram

import (
	"encoding/json"
	"testing"
)

func TestTokenResponse_NumericUserID(t *testing.T) {
	payload := `{"access_token": "IGAA2e0uoy3ZAZABZA...", "user_id": 27640875392278831, "permissions": ["instagram_business_basic"]}`

	var tokenResp TokenResponse
	if err := json.Unmarshal([]byte(payload), &tokenResp); err != nil {
		t.Fatalf("Failed to unmarshal TokenResponse with numeric user_id: %v", err)
	}

	if tokenResp.AccessToken != "IGAA2e0uoy3ZAZABZA..." {
		t.Errorf("Expected token, got %s", tokenResp.AccessToken)
	}

	if tokenResp.GetUserID() != "27640875392278831" {
		t.Errorf("Expected user_id string '27640875392278831', got %s", tokenResp.GetUserID())
	}
}

func TestTokenResponse_StringUserID(t *testing.T) {
	payload := `{"access_token": "IGAA2e0uoy3ZAZABZA...", "user_id": "27640875392278831"}`

	var tokenResp TokenResponse
	if err := json.Unmarshal([]byte(payload), &tokenResp); err != nil {
		t.Fatalf("Failed to unmarshal TokenResponse with string user_id: %v", err)
	}

	if tokenResp.GetUserID() != "27640875392278831" {
		t.Errorf("Expected user_id string '27640875392278831', got %s", tokenResp.GetUserID())
	}
}
