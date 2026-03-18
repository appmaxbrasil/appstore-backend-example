package controllers_test

import (
	"os"
	"testing"

	"github.com/appmaxbrasil/appstore-backend-example/tests/unit/support"
)

func TestMain(m *testing.M) {
	os.Exit(support.RunWithFrameworkBootstrap(m))
}
