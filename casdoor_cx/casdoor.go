package casdoor_cx

import (
	"context"
	"fmt"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"golang.org/x/oauth2"
	"os"
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

func NewCasDoorClientConfig(clientConfig *CasDoorClientConfig) (*casdoorsdk.Client, error) {
	// 加载一下本地的证书文件
	if clientConfig.CertificateName != "" && clientConfig.Certificate == "" {
		raw, err := os.ReadFile("./certs/" + clientConfig.CertificateName + ".pem")
		if err != nil {
			return nil, err
		}
		if len(raw) > 0 {
			clientConfig.Certificate = string(raw)
		}
	}

	client := casdoorsdk.NewClient(clientConfig.Endpoint, clientConfig.ClientId, clientConfig.ClientSecret, clientConfig.Certificate,
		clientConfig.OrganizationName, clientConfig.ApplicationName)

	// 远程的证书
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

func NewCasDoorClient(clientConfig *CasDoorClientConfig) (*CasDoorClient, error) {
	client, err := NewCasDoorClientConfig(clientConfig)
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

type Claims *casdoorsdk.Claims

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
				claimB := Claims(claim) // 类型转换一下
				ctx = context.WithValue(ctx, AuthKey{}, claimB)
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
