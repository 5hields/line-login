package line_login

import (
	"github.com/5hields/line-login/linethrift"
	"context"
	"log"
)

var ctx = context.Background()

func (c *LineClient) GetRSAKeyInfo() *linethrift.RSAKey {
	res, err := c.talkClient.GetRSAKeyInfo(ctx, linethrift.IdentityProvider_LINE)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func (c *LineClient) RequestAccountPasswordReset(email string) {
	if err := c.talkClient.RequestAccountPasswordReset(ctx, linethrift.IdentityProvider_LINE, email, "ja"); err != nil {
		log.Fatal(err)
	}

}

func (c *LineClient) ConfirmE2EELogin(verifier string, secret []byte) string {
	res, err := c.authClient.ConfirmE2EELogin(ctx, verifier, secret)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func (c *LineClient) LoginZ(req *linethrift.LoginRequest) *linethrift.LoginResult_ {
	res, err := c.authClient.LoginZ(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func (c *LineClient) LogoutZ() {
	if err := c.authClient.LogoutZ(ctx); err != nil {
		log.Fatal(err)
	}
}
