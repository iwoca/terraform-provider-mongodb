package mongodb

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"mongodb": func() (*schema.Provider, error) { //nolint:unparam // signature required by resource.TestCase.ProviderFactories
			return Provider(), nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("MONGODB_HOST"); v == "" {
		t.Fatal("MONGODB_HOST must be set for acceptance tests")
	}
	if v := os.Getenv("MONGODB_USERNAME"); v == "" {
		t.Fatal("MONGODB_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("MONGODB_PASSWORD"); v == "" {
		t.Fatal("MONGODB_PASSWORD must be set for acceptance tests")
	}
}
