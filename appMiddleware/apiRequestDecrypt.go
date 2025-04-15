package appMiddleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/doubunv/common-pkg/aesGCM"
	"github.com/doubunv/common-pkg/result"
	"io"
	"net/http"
)

var RequestDecryptError = errors.New("Request decryption failed. ")

type ApiRequestDecryptOption func(m *ApiRequestDecryptMiddleware)

type RequestDecryptData struct {
	AesData string `json:"aes_data"`
}

type ApiRequestDecryptMiddleware struct {
	aesKey string
	debug  bool
}

func DecryptKeyOption(aesKey string) ApiRequestDecryptOption {
	return func(m *ApiRequestDecryptMiddleware) {
		m.aesKey = aesKey
	}
}

func DecryptWithDebugOption() ApiRequestDecryptOption {
	return func(m *ApiRequestDecryptMiddleware) {
		m.debug = true
	}
}

func NewApiRequestDecryptMiddleware(arg ...ApiRequestDecryptOption) *ApiRequestDecryptMiddleware {
	res := &ApiRequestDecryptMiddleware{}
	for _, o := range arg {
		o(res)
	}

	return res
}

func (m *ApiRequestDecryptMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if aesGCM.IsOpenAesGcm {
			if err := m.RequestDecrypt(r); err != nil {
				result.HttpErrorResult(r.Context(), w, err)
				return
			}
		}

		next(w, r)
	}
}

func (m *ApiRequestDecryptMiddleware) RequestDecrypt(r *http.Request) error {
	if m.debug {
		return nil
	}
	aesGCM.EncryptKey = []byte(m.aesKey)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return RequestDecryptError
	}

	if len(data) == 0 {
		return nil
	}
	// Decrypt the data here
	var decryptData RequestDecryptData
	if err = json.Unmarshal(data, &decryptData); err != nil {
		return RequestDecryptError
	}

	deData, err := aesGCM.Decrypt(aesGCM.EncryptKey, decryptData.AesData)
	if err != nil {
		return RequestDecryptError
	}

	r.Body = io.NopCloser(bytes.NewBuffer(deData))

	return nil
}
