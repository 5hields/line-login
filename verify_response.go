package line_login

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

type verifyResponse struct {
	Result    result `json:"result"`
	Timestamp string `json:"timestamp"`
	AuthPhase string `json:"authPhase"`
}

type result struct {
	Metadata struct {
		ErrorCode         string `json:"errorCode"`
		EncryptedKeyChain string `json:"encryptedKeyChain"`
		E2eeVersion       string `json:"e2eeVersion"`
		KeyId             string `json:"keyId"`
		PublicKey         string `json:"publicKey"`
	} `json:"metadata"`
}

func waitingVerifier(verifier string) ([]byte, []byte) {
	res, err := getVerifyURL(verifier)
	if err != nil {
		log.Fatal(err)
	}

	verifyRes, err := decodeToVerifyResponse(res)
	if err != nil {
		log.Fatal(err)
	}

	if verifyRes.IsCanceled() {
		log.Fatal("Login Process canceled.")
	}

	pubKey, keyChain := verifyRes.GetE2EEKeys()
	return pubKey, keyChain
}

func decodeToVerifyResponse(res *http.Response) (*verifyResponse, error) {
	decoder := json.NewDecoder(res.Body)
	result := verifyResponse{}
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return &result, nil
}

func (r *verifyResponse) GetE2EEKeys() ([]byte, []byte) {
	meta := r.Result.Metadata
	decPubKey, _ := base64.StdEncoding.DecodeString(meta.PublicKey)
	decKeyChain, _ := base64.StdEncoding.DecodeString(meta.EncryptedKeyChain)

	return decPubKey, decKeyChain
}

func (r *verifyResponse) IsCanceled() bool {
	if r.Result.Metadata.ErrorCode == "CANCEL" {
		return true
	}
	return false
}

func getVerifyURL(verifier string) (*http.Response, error) {
	req, err := http.NewRequest("GET", verifyUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Line-Access", verifier)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Line-Application", lineApp)

	hc := http.Client{}
	resp, _ := hc.Do(req)

	return resp, nil
}
