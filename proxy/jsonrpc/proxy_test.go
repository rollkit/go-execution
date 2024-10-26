package jsonrpc_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	jsonrpcproxy "github.com/LastL2/go-execution/proxy/jsonrpc"
	"github.com/LastL2/go-execution/test"
)

type ProxyTestSuite struct {
	test.ExecuteSuite
	server  *httptest.Server
	client  *jsonrpcproxy.Client
	cleanup func()
}

func (s *ProxyTestSuite) SetupTest() {
	exec := test.NewExecute()
	config := &jsonrpcproxy.Config{
		DefaultTimeout: time.Second,
		MaxRequestSize: 1024 * 1024,
	}
	server := jsonrpcproxy.NewServer(exec, config)

	s.server = httptest.NewServer(server)

	client := jsonrpcproxy.NewClient()
	client.SetConfig(config)

	err := client.Start(s.server.URL)
	require.NoError(s.T(), err)

	s.client = client
	s.Exec = client
	s.cleanup = func() {
		client.Stop()
		s.server.Close()
	}
}

func (s *ProxyTestSuite) TearDownTest() {
	s.cleanup()
}

func TestProxySuite(t *testing.T) {
	suite.Run(t, new(ProxyTestSuite))
}
