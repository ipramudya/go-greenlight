package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl), // provide ttl(time to live) duration parameter to the current expiry time
		Scope:  scope,
	}

	// initiate zero-valued slice with length of 16 bytes
	randomBytes := make([]byte, 16)

	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}

	// Encode the byte slice to a base-32-encoded string and assign it to the token
	// This will be the token string that we send to the user in ther Welcome email.
	// They will look similar to "Y3QMGX3PJ3WLRL2YRTQGQ6KRHU"
	// Note that by default base-32 strings may be padded at the end with the = character
	// We don't need this padding character for the purpose of our tokens, so we use
	// the WithPadding(base32.NoPadding) method in the line below to omit them
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Generate a SHA-256 hash of the plaintext token string. This will be the value
	// that we store in the `hash` field of our database table. Note that the
	// sha256.Sum256() function returns an *array* of length 32, so to make it easier to
	// work with we convert it to a slice using the [:] operator before storing it.
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}
