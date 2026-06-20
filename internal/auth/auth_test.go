package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPaswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}
func TestCheckBearerToken(t *testing.T) {
	t.Log("Hi")
	b := []string{"quasi", "hello", "sup", "yes", "no"}
	headers := make([]http.Header, 5)
	for i := 0; i < 5; i++ {
		a := http.Header{}
		a.Add("Authorization", fmt.Sprintf("Bearer %s", b[i]))
		headers[i] = a
	}
	for i := 0; i < 5; i++ {
		tok, err := GetBearerToken(headers[i])
		t.Log(tok)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestCheckJWT(t *testing.T) {
	testUUID := uuid.New()

	tokenSecret := "my-secret-token"
	dur, err := time.ParseDuration("20s")
	if err != nil {
		t.Fatal(err)
	}
	jwt, err := MakeJWT(testUUID, tokenSecret, dur)
	if err != nil {
		t.Fatal(err)
	}
	validatedJWT, err := ValidateJWT(jwt, tokenSecret)
	if err != nil {
		t.FailNow()
	}
	if validatedJWT != testUUID {
		t.Fatal(err)
	}

}
func TestCheckJWTTimeout(t *testing.T) {
	testUUID := uuid.New()

	tokenSecret := "my-secret-token"
	jwt, err := MakeJWT(testUUID, tokenSecret, -1*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateJWT(jwt, tokenSecret)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

}

func TestWrongSecret(t *testing.T) {
	testUUID := uuid.New()

	tokenSecret := "my-secret-token"
	dur, err := time.ParseDuration("20s")
	if err != nil {
		t.Fatal(err)
	}
	jwt, err := MakeJWT(testUUID, tokenSecret, dur)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ValidateJWT(jwt, "wrong-secret")
	if err == nil {
		t.Fatal("expected error but got nil")
	}

}
