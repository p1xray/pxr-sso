package cookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/p1xray/pxr-sso/internal/oidc/domain/dto"
	"io"
	"time"
)

type SessionEncoding interface {
	Encode(session dto.AuthorizedSession) (string, error)
	Decode(encoded string) (dto.AuthorizedSession, error)
}

type Session struct {
	SessionID        string    `json:"sid"`
	Subject          int64     `json:"sub"`
	AuthTime         time.Time `json:"auth_time"`
	IdentityProvider string    `json:"identity_provider"`
}

func NewSession(code string, userID int64, authTime time.Time) Session {
	return Session{
		SessionID:        code,
		Subject:          userID,
		AuthTime:         authTime,
		IdentityProvider: "pxr.sso",
	}
}

type sessionEncoding struct {
	key []byte
}

func NewSessionEncoding(cfg SessionConfig) *sessionEncoding {
	return &sessionEncoding{
		key: []byte(cfg.SecretKey),
	}
}

func (se *sessionEncoding) Encode(session dto.AuthorizedSession) (string, error) {
	sessionToMarshal := Session{
		SessionID:        session.ID(),
		Subject:          session.Subject(),
		AuthTime:         session.AuthTime(),
		IdentityProvider: session.IdentityProvider(),
	}

	sessionJSON, err := json.Marshal(sessionToMarshal)
	if err != nil {
		return "", fmt.Errorf("marshalling session: %w", err)
	}

	encryptedSession, err := se.encrypt(sessionJSON)
	if err != nil {
		return "", fmt.Errorf("encrypting session: %w", err)
	}

	return encryptedSession, nil
}

func (se *sessionEncoding) Decode(encoded string) (dto.AuthorizedSession, error) {
	decryptedSession, err := se.decrypt(encoded)
	if err != nil {
		return dto.AuthorizedSession{}, fmt.Errorf("decrypting session: %w", err)
	}

	var session Session
	err = json.Unmarshal(decryptedSession, &session)
	if err != nil {
		return dto.AuthorizedSession{}, fmt.Errorf("unmarshalling session: %w", err)
	}

	authorizedSession := dto.NewAuthorizedSession(
		session.SessionID,
		session.Subject,
		session.AuthTime,
		session.IdentityProvider,
	)

	return authorizedSession, nil
}

func (se *sessionEncoding) encrypt(value []byte) (string, error) {
	block, err := aes.NewCipher(se.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, value, nil)
	return base64.RawURLEncoding.EncodeToString(cipherText), nil
}

func (se *sessionEncoding) decrypt(value string) ([]byte, error) {
	cipherText, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(se.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()

	nonce, actualCipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainTextBytes, err := aesGCM.Open(nil, nonce, actualCipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainTextBytes, nil
}
