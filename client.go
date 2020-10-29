package line_login

import (
	"github.com/5hields/line-login/linethrift"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
)

func newThriftClient(apiUrl string) (*thrift.TStandardClient, error) {
	trans, err := thrift.NewTHttpClient(apiUrl)
	if err != nil {
		return nil, err
	}
	httpTrans := trans.(*thrift.THttpClient)
	header := map[string]string{
		"X-Line-Application": lineApp,
		"User-Agent":         userAgent,
	}
	for key, val := range header {
		httpTrans.SetHeader(key, val)
	}
	protocol := thrift.NewTCompactProtocolFactory().GetProtocol(trans)
	thriftClient := thrift.NewTStandardClient(protocol, protocol)

	return thriftClient, nil
}

type LineClient struct {
	talkClient *linethrift.TalkServiceClient
	authClient *linethrift.AuthServiceClient
}

func NewLineClient() *LineClient {
	talk, err := newThriftClient(registerUrl)
	if err != nil {
		log.Fatal(err)
	}
	auth, err := newThriftClient(authRegisterUrl)
	if err != nil {
		log.Fatal(err)
	}
	return &LineClient{
		talkClient: linethrift.NewTalkServiceClient(talk),
		authClient: linethrift.NewAuthServiceClient(auth),
	}
}
