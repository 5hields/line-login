package line_login

import (
	"github.com/5hields/line-login/linethrift"
	"fmt"
)

type loginOption struct {
	Cert string
}

type option func(*loginOption)

func newLoginOption() *loginOption {
	return &loginOption{
		Cert: "",
	}
}

func Certificate(cert string) func(*loginOption) {
	return func(l *loginOption) {
		l.Cert = cert
	}
}

func LoginResultDetail(r *linethrift.LoginResult_) {
	fmt.Println("Logging in Successful")
	fmt.Println("Certificate:", r.Certificate)
	fmt.Println("AuthToken:", r.AuthToken)
}

func newLoginRequestWithE2EE(key *linethrift.RSAKey, secret []byte, encryptedPwd, cert string) *linethrift.LoginRequest {
	return &linethrift.LoginRequest{
		Type:             linethrift.LoginType_ID_CREDENTIAL_WITH_E2EE,
		IdentityProvider: linethrift.IdentityProvider_LINE,
		Identifier:       key.Keynm,
		Password:         encryptedPwd,
		KeepLoggedIn:     true,
		Certificate:      cert,
		SystemName:       "Chino",
		Secret:           secret,
		E2eeVersion:      1,
	}
}

func updateLoginRequest(r *linethrift.LoginRequest, verifier string) *linethrift.LoginRequest {
	r.Type = linethrift.LoginType_QRCODE
	r.Verifier = verifier
	return r
}

func LoginWithCredential(email, pass string, opts ...func(*loginOption)) (string, error) {
	line := NewLineClient()
	rsaKey := line.GetRSAKeyInfo()
	cred := encryptCredential(email, pass, rsaKey)

	keyPair := newE2EEKeyPair()
	pin, secret := createSecret(keyPair.Pub[:])

	c := newLoginOption()
	for _, opt := range opts {
		opt(c)
	}

	req := newLoginRequestWithE2EE(rsaKey, secret, cred, c.Cert)
	res := line.LoginZ(req)

	switch res.Type {
	case linethrift.LoginResultType_SUCCESS:
		res.Certificate = c.Cert
		LoginResultDetail(res)
		return res.AuthToken, nil
	case linethrift.LoginResultType_REQUIRE_DEVICE_CONFIRM:
		fmt.Println("Enter [", pin, "] on your mobile device")
		pubKey, keyChain := waitingVerifier(res.Verifier)
		deviceSecret := getHashKeyChain(keyPair.GenSharedSecret(pubKey), keyChain)
		verifier := line.ConfirmE2EELogin(res.Verifier, deviceSecret)

		/* LoginZ E2EE DeviceConfirm */
		updateLoginRequest(req, verifier)
		lRes := line.LoginZ(req)
		LoginResultDetail(lRes)
		return lRes.AuthToken, nil
	default:
		return "", fmt.Errorf("SecondaryLogin Failed.\n")
	}
}
