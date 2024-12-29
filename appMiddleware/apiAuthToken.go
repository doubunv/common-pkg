package appMiddleware

import (
	"github.com/doubunv/common-pkg/consts"
	"github.com/doubunv/common-pkg/ctxMd"
	"github.com/doubunv/common-pkg/headInfo"
	"github.com/doubunv/common-pkg/result/xcode"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strconv"
	"strings"
)

type CheckRequestTokenFunc func(r *http.Request, token string) (userId int64, roleId int64)

func verifyPath(urlPath string, noVerifyPath map[string]int) bool {
	if _, ok := noVerifyPath[urlPath]; ok {
		return true
	}

	for path, _ := range noVerifyPath {
		if strings.HasPrefix(path, "/") && strings.HasSuffix(path, "*") {
			prefix := strings.TrimSuffix(path, "*")
			if strings.HasPrefix(urlPath, prefix) {
				return true
			}
		}
	}
	return false
}

func MustAuthTokenRequest(r *http.Request, checkToken CheckRequestTokenFunc, noVerifyPath map[string]int) (*http.Request, error) {
	ctx := r.Context()
	token := headInfo.GetJwtToken(ctx)

	wPathBool := verifyPath(r.URL.Path, noVerifyPath)
	if !wPathBool && token == "" {
		return r, xcode.UserNotFound
	}

	if token != "" {
		var tokenUid int64 = 0
		var roleId int64 = 0
		if checkToken != nil {
			tokenUid, roleId = checkToken(r, token)
		}
		if tokenUid > 0 {
			md := ctxMd.SetMdCtxFromOut(ctx, consts.TokenUid, strconv.FormatInt(tokenUid, 10))
			md = ctxMd.SetMdCtxFromOut(ctx, consts.TokenUidRole, strconv.FormatInt(roleId, 10))
			ctx = metadata.NewOutgoingContext(ctx, md)
		} else {
			if !wPathBool {
				return r, xcode.TokenInvalid
			}
		}
	}

	newReq := r.WithContext(ctx)

	return newReq, nil
}
