package buffalo

import (
	"github.com/gavv/httpexpect"
	"github.com/wowucco/go-admin/tests/common"
	"net/http"
	"testing"
)

func TestBuffalo(t *testing.T) {
	common.ExtraTest(httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(newHandler()),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	}))
}
