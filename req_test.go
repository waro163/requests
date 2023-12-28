package requests

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoBaseUrl(t *testing.T) {
	url := "https://baidu.com"
	client := NewClient()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)
	resp, err := client.DoRequest(req.Context(), req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	t.Logf("%s", body)
}
