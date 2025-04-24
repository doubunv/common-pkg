package appMiddleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/doubunv/common-pkg/aesGCM"
	"github.com/doubunv/common-pkg/consts"
	"github.com/doubunv/common-pkg/result"
	"io"
	"net/http"
)

var RequestDecryptError = errors.New("Request decryption failed. ")

var RequestBadError = errors.New("Request bad. ")

type ApiRequestDecryptOption func(m *ApiRequestDecryptMiddleware)

type RequestDecryptData struct {
	AesData string `json:"aes_data,default=''"`
}

type ApiRequestDecryptMiddleware struct {
}

func DecryptKeyOption(aesKey string) ApiRequestDecryptOption {
	return func(m *ApiRequestDecryptMiddleware) {
		aesGCM.EncryptKey = []byte(aesKey)
	}
}

func DecryptWithDebugFalseOption() ApiRequestDecryptOption {
	return func(m *ApiRequestDecryptMiddleware) {
		aesGCM.IsOpenAesGcm = true
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
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	if len(data) == 0 {
		return nil
	}
	// Decrypt the data here
	var decryptData RequestDecryptData
	if err = json.Unmarshal(data, &decryptData); err != nil {
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		return nil
	}

	if decryptData.AesData == "" && r.Header.Get(consts.BusinessCode) != "" {
		return RequestBadError
	}

	if decryptData.AesData == "" {
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		return nil
	}

	deData, err := aesGCM.Decrypt(aesGCM.EncryptKey, decryptData.AesData)
	if err != nil {
		return RequestDecryptError
	}

	if deData == nil {
		deData = []byte("{}")
	}
	r.Body = io.NopCloser(bytes.NewBuffer(deData))
	return nil
}
