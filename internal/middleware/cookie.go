package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type CookieHandler struct {
	aesgcm cipher.AEAD
	nonce  []byte
}

var (
	ch                    *CookieHandler
	cookieUserIDFieldName = "user_id"
	secretKey             = "Do you see the gopher? No. Me neither. But he is exist."
)

func NewCookieHandler() error {
	if ch == nil {
		key := sha256.Sum256([]byte(secretKey))
		aesblock, err := aes.NewCipher(key[:])
		if err != nil {
			return err
		}
		aesgcm, err := cipher.NewGCM(aesblock)
		if err != nil {
			return err
		}
		nonce := key[len(key)-aesgcm.NonceSize():]

		ch = &CookieHandler{
			aesgcm: aesgcm,
			nonce:  nonce,
		}
	}

	return nil
}

func UserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := uuid.New().String()

		err := GetCurrentUserID(r, &userID)

		if err != nil && err != http.ErrNoCookie {
			io.WriteString(w, err.Error())
			return
		}

		enc, err := Encode(userID)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		cookie := &http.Cookie{
			Name:  cookieUserIDFieldName,
			Value: enc,
			Path:  "/",
		}

		http.SetCookie(w, cookie)
		r.AddCookie(cookie)

		next.ServeHTTP(w, r)
	})
}

func GetCurrentUserID(r *http.Request, userID *string) error {
	cookieUserID, err := r.Cookie(cookieUserIDFieldName)
	if err != nil {
		return err
	}

	err = Decode(cookieUserID.Value, userID)
	if err != nil {
		return err
	}

	return nil
}

func Decode(encUserID string, userID *string) error {
	err := NewCookieHandler()
	dst, err := hex.DecodeString(encUserID)
	if err != nil {
		return err
	}

	src, err := ch.aesgcm.Open(nil, ch.nonce, dst, nil)
	if err != nil {
		return err
	}
	*userID = string(src)

	return nil
}

func Encode(userID string) (string, error) {
	err := NewCookieHandler()
	if err != nil {
		return "", nil
	}

	src := []byte(userID)
	enc := ch.aesgcm.Seal(nil, ch.nonce, src, nil)
	return hex.EncodeToString(enc), nil
}
