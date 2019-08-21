package client

import (
	"github.com/Bplotka/go-httpt"
	"github.com/Bplotka/go-httpt/rt"
	"net/http"
	"testing"
)

func TestGetDeploymentWithMockedHttpClient(t *testing.T) {

	// given
	s := mockLiimaServer(t)
	testClient := s.HTTPClient()

	cli := Cli{}
	cli.Client = NewMockClientWithCustomHttpClient(testClient)

	// when
	testResponse := TestResponse{}
	if err := cli.Client.DoRequest(http.MethodGet, "deployments/test", nil, &testResponse); err != nil {
		t.Errorf("Excepting no error: %s", err)
	}

	// then
	assertString(t, "test", testResponse.Name, "name")
	assertString(t, "252734", testResponse.Value, "value")
}

func assertString(t *testing.T, expected string, result string, fieldName string) {
	if result != expected {
		t.Errorf("%s doesn't seem correct, got %s but expected %s.", fieldName, result, expected)
	}
}
func mockLiimaServer(t *testing.T) *httpt.Server {
	s := httpt.NewServer(t)
	s.On(httpt.GET, "deployments/test").Push(rt.JSONResponseFunc(http.StatusOK, []byte(
		`{
        "name": "test",
		"value": "252734"
	 	}`,
	)))
	return s
}

type TestResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
