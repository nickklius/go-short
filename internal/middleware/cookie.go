package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type CookieHandler struct {
	aesgcm cipher.AEAD
	nonce  []byte
}

var (
	ch                    CookieHandler
	once                  sync.Once
	cookieInitErr         error
	cookieUserIDFieldName = "user_id"
	secretKey             = "Do you see the gopher? No. Me neither. But he is exist."
)

func NewCookieHandler() (CookieHandler, error) {
	once.Do(func() {
		key := sha256.Sum256([]byte(secretKey))

		var aesblock cipher.Block
		aesblock, cookieInitErr = aes.NewCipher(key[:])
		if cookieInitErr != nil {
			return
		}

		var aesgcm cipher.AEAD
		aesgcm, cookieInitErr = cipher.NewGCM(aesblock)
		if cookieInitErr != nil {
			return
		}
		nonce := key[len(key)-aesgcm.NonceSize():]

		ch = CookieHandler{
			aesgcm: aesgcm,
			nonce:  nonce,
		}
	})

	return ch, cookieInitErr
}

func UserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetCurrentUserID(r)

		if err != nil && err != http.ErrNoCookie {
			io.WriteString(w, err.Error())
			return
		}

		if userID == "" {
			userID = uuid.New().String()
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

func GetCurrentUserID(r *http.Request) (string, error) {
	cookieUserID, err := r.Cookie(cookieUserIDFieldName)
	if err != nil {
		return "", err
	}

	var userID string
	err = Decode(cookieUserID.Value, &userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func Decode(encUserID string, userID *string) error {
	c, err := NewCookieHandler()
	if err != nil {
		return err
	}

	dst, err := hex.DecodeString(encUserID)
	if err != nil {
		return err
	}

	src, err := c.aesgcm.Open(nil, c.nonce, dst, nil)
	if err != nil {
		return err
	}
	*userID = string(src)

	return nil
}

func Encode(userID string) (string, error) {
	c, err := NewCookieHandler()
	if err != nil {
		return "", nil
	}

	src := []byte(userID)
	enc := c.aesgcm.Seal(nil, c.nonce, src, nil)
	return hex.EncodeToString(enc), nil
}
