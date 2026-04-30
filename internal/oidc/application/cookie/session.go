package cookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
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

	return base64.StdEncoding.EncodeToString(encryptedSession), nil
}

func (se *sessionEncoding) Decode(encoded string) (dto.AuthorizedSession, error) {
	encodedAsBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return dto.AuthorizedSession{}, fmt.Errorf("decoding session: %w", err)
	}

	decryptedSession, err := se.decrypt(encodedAsBytes)
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

func (se *sessionEncoding) encrypt(value []byte) ([]byte, error) {
	block, err := aes.NewCipher(se.key)
	if err != nil {
		return nil, err
	}

	b := base64.StdEncoding.EncodeToString(value)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCTR(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	return ciphertext, nil
}

func (se *sessionEncoding) decrypt(value []byte) ([]byte, error) {
	block, err := aes.NewCipher(se.key)
	if err != nil {
		return nil, err
	}

	if len(value) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := value[:aes.BlockSize]
	value = value[aes.BlockSize:]
	cfb := cipher.NewCTR(block, iv)
	cfb.XORKeyStream(value, value)
	data, err := base64.StdEncoding.DecodeString(string(value))
	if err != nil {
		return nil, err
	}

	return data, nil
}
