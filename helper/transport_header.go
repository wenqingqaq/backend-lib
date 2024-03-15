package helper

import (
	"context"
	"encoding/json"
	errno "gitee.com/yanwenqing/backend-lib/api/err/v1"

	"gitee.com/yanwenqing/backend-lib/casdoor"
	"gitee.com/yanwenqing/backend-lib/tecons"
	"github.com/go-kratos/kratos/v2/transport"
)

func UserFromHeader(ctx context.Context) (*casdoor.User, error) {
	c, err := ClaimsFromHeader(ctx)
	if err != nil {
		return nil, err
	}
	return &c.User, nil
}

func ClaimsFromHeader(ctx context.Context) (*casdoor.Claims, error) {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		//return nil, kjwt.ErrMissingJwtToken
		return nil, errno.ErrorAuthMissingJwtToken("")
	}
	var c casdoor.Claims
	jsonStr := tr.RequestHeader().Get(tecons.TransHeaderClaims)
	if jsonStr == "" {
		return nil, errno.ErrorAuthMissingTecoClaims("ClaimsFromHeader:Teco_claims not found")
	}
	if err := json.Unmarshal([]byte(jsonStr), &c); err != nil {
		return nil, errno.ErrorUnknown("ClaimsFromHeader,Teco_claims:%s,err:%s", jsonStr, err.Error())
	}
	return &c, nil
}

func TenantIDFromHeader(ctx context.Context) (string, error) {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		//return "", er.ErrorTenantIdInvalid("")
		return "", errno.ErrorAuthMissingJwtToken("")
	}
	cID := tr.RequestHeader().Get(tecons.TransHeaderTenantID)
	if cID == "" {
		return cID, errno.ErrorAuthInvalidTenantId("")
	}
	return cID, nil
}
