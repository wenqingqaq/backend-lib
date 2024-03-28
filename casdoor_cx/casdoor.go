package casdoor_cx

import (
	"context"
	"fmt"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"golang.org/x/oauth2"
	"strings"
)

type CasDoorClientConfig struct {
	Endpoint         string
	ClientId         string
	ClientSecret     string
	Certificate      string
	CertificateName  string // 证书名字
	OrganizationName string
	ApplicationName  string
}

func NewCasDoorClientConfig() (*casdoorsdk.Client, error) {
	clientConfig := CasDoorClientConfig{
		Endpoint:         "http://localhost:8100",
		ClientId:         "6be1ec496f3173e35509",
		ClientSecret:     "a587bba2075f864dc393a07ace30524f246159f6",
		Certificate:      "",
		OrganizationName: "built-in",
		ApplicationName:  "云平台",
		CertificateName:  "cert-built-in",
	}

	// 加载一下本地的证书文件
	//raw, err := os.ReadFile("./certs/" + clientConfig.CertificateName + ".pem")
	//if err != nil {
	//	return nil, err
	//}
	//if len(raw) > 0 {
	//	clientConfig.Certificate = string(raw)
	//}

	clientConfig.Certificate = `-----BEGIN CERTIFICATE-----
MIIE3TCCAsWgAwIBAgIDAeJAMA0GCSqGSIb3DQEBCwUAMCgxDjAMBgNVBAoTBWFk
bWluMRYwFAYDVQQDEw1jZXJ0LWJ1aWx0LWluMB4XDTI0MDMyNzAxMzI0OFoXDTQ0
MDMyNzAxMzI0OFowKDEOMAwGA1UEChMFYWRtaW4xFjAUBgNVBAMTDWNlcnQtYnVp
bHQtaW4wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDxwskU0x4nGcp9
w+A4iVzZNx6Imv41HmlGewF9o5Kt7/InhZV3JFPlekVBaoGLfRpbPuVo3gEKevPW
T9Fc1AyKcy2Q7+rERtxsY/THhoyK9KAfm6WO30RTc5lAJA4IHQG/0YJTwC1I8Hwi
Ar3xeN3zqMt5bG5ifcXY2NpIJMRtKNBrSaaJkeIMn/fAUT5DtrwxOKdc8BfVgKOm
40lOtwwxkCjeyDN6chEnpU+NOoskXGjdlsSBA3vz+DB3vOppBlRPiars59loVZMc
ayCpSpe4tkt8dnThmQtmD37TE3I0RdWW40FG3R7AfovIqH+0Hi7OYtyRX2zBXIB2
UOk8CYDtjImAW5SM0CovKLiJE85Q+wscZBbDN/fV7AujHl2UXNcICgEnGOXT1/SH
i7hUqM7x7BzWRDu6bNFVExSJPNlJ3yzqXfbFnCKrk55eb9hwYe4ZO4/iR2qV0Ivw
tqECIBgLQcFYBWVGXUPH046aedTg+jIepebZLzmw6ETXeWgHPNNfE1yeL6vadNrn
TLgRCgaoyxhvJSkxyWTV9GDI42sYV5ZOxZE7LZcXCLBIRSo0a8OyIsQAJDvPZUee
jReC6bS221GFL1GkLmkEWuZk2IeQSmFObiunaqG3wFVFPL4iukiqrfiFmR2Kuyf2
XWQS2Eq2yW0z6n88HoKt1wXLPObicwIDAQABoxAwDjAMBgNVHRMBAf8EAjAAMA0G
CSqGSIb3DQEBCwUAA4ICAQDvejeCdADX6kyFf19kQgdMharQsNGbYnlv/5QlMV/l
u9m2T+ODCQkehJKgA3idUOqdOxWuvbS5720PJRgM2V4OMKly2ZQYPqGPb1yMGna5
96mZwXhaIDlRBtOO7mFfe0xpNUlN00ucDvkybIiCp8b9GABeGaCSLMfDwZjcxA8F
ucrGY+oMuBsRm0uG8gVUMGRkhwYq5CjSUVTtuB4mo5z+csA96FFjyMMngUppW7sh
wBBfLpJvvrB0rvm60o0mRLTbSuTVodcCcGaO8WHgJOWwpe797nF7/eWaSGv3S0tu
AgCYgBtC2pRQZzx0pdwZ8swtUzOnMZZ7QX1L5VKfG3rYgWj0C29weRBSBD/6A27D
xiS6qiwf4DUD+7L2m1Ckf6yKwiJkCfz2RVIzBS5zb/+5h+WrPkHqXKTQKIKQgCEQ
hAzOaUNWjrz2acvmkhKyfvNI0exChce+bSwOOvj6QEowBYhLMYSUw+WBYQSLgQxQ
YzsNjPIZF+UgwRVcxA7OwgwB/vEOkoba8SS3pZLQ587jQEnNWcKdaGmjblLOB2Ex
7vfd1goNj5xrKaz10OFSrSvtKvdnGR3pof6Dzu7VD5qkf+8gg8mENriPEbK+mm1I
d+9A3Nmvf/NKu62pVcX8k8CqPuxmKMXmhtJWh4b9DFXdHr1bgo5hxKBEVLTuMI6q
uA==
-----END CERTIFICATE-----`

	client := casdoorsdk.NewClient(clientConfig.Endpoint, clientConfig.ClientId, clientConfig.ClientSecret, clientConfig.Certificate,
		clientConfig.OrganizationName, clientConfig.ApplicationName)

	if client.Certificate == "" {
		certs, err := client.GetCerts()
		if err != nil {
			return nil, err
		}
		for _, cert := range certs {
			if cert.Name == clientConfig.CertificateName {
				client.Certificate = cert.Certificate
			}
		}
		if client.Certificate == "" {
			return nil, errors.New(200, "certificate not found", "certificate not found")
		}
	}

	return client, nil
}

type CasDoorClient struct {
	Client *casdoorsdk.Client
}

func NewCasDoorClient() (*CasDoorClient, error) {
	client, err := NewCasDoorClientConfig()
	if err != nil {
		return nil, err
	}

	return &CasDoorClient{Client: client}, nil
}

func (c *CasDoorClient) GetOAuthToken(code, state string) *oauth2.Token {
	token, err := c.Client.GetOAuthToken(code, state)
	if err != nil {
		panic(err)
	}
	fmt.Println(token.AccessToken)

	return token
}

const (

	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// bearerFormat authorization token format
	bearerFormat string = "Bearer %s"

	// authorizationKey holds the key used to store the JWT Token in the request tokenHeader.
	authorizationKey string = "Authorization"

	// reason holds the error reason.
	reason string = "UNAUTHORIZED"
)

var (
	ErrMissingJwtToken = errors.Unauthorized(reason, "JWT token is missing")
	ErrTokenInvalid    = errors.Unauthorized(reason, "Token is invalid")
	ErrTokenInvalid2   = errors.Unauthorized(reason, "22222 is invalid")
	ErrTokenClaim      = errors.Unauthorized(reason, "Token claim error")
	ErrWrongContext    = errors.Unauthorized(reason, "Wrong context for middleware")
)

type AuthKey struct{}

type Claims casdoorsdk.Claims

func (c *CasDoorClient) CasDoorJWT() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				auths := strings.SplitN(header.RequestHeader().Get(authorizationKey), " ", 2)
				if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
					return nil, ErrMissingJwtToken
				}
				jwtToken := auths[1]
				if jwtToken == "" {
					return nil, ErrTokenInvalid
				}
				claim, err := c.Client.ParseJwtToken(jwtToken)
				if err != nil {
					return nil, err
				}
				ctx = context.WithValue(ctx, "user", claim)
				return handler(ctx, req)
			} else {
				return nil, ErrTokenInvalid2
			}
		}
	}
}

func (c *CasDoorClient) CustomerJWT() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			claims, ok := ctx.Value(AuthKey{}).(Claims)
			fmt.Println(claims, ok)
			if !ok {
				// 没有获取到 claims
				return nil, errors.Unauthorized("UNAUTHORIZED", "claims not found")
			}

			ctx = context.WithValue(ctx, "user_id", claims.User.Id)

			// 可以查询一下用户权限

			return handler(ctx, req)
		}
	}
}
