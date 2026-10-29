package mongodb

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMongoDBUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBUserConfig("read"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mongodb_db_user.test", "name", "testuser"),
					resource.TestCheckResourceAttr("mongodb_db_user.test", "auth_database", "admin"),
				),
			},
			{
				ResourceName:            "mongodb_db_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccMongoDBUser_updateRoles(t *testing.T) {
	// Test that updating roles doesn't drop the user (in-place update)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBUserConfig("read"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mongodb_db_user.test", "name", "testuser"),
				),
			},
			{
				Config: testAccMongoDBUserConfig("readWrite"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mongodb_db_user.test", "name", "testuser"),
				),
			},
		},
	})
}

func TestAccMongoDBUser_updatePassword(t *testing.T) {
	// Test that updating password doesn't drop the user (in-place update)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBUserConfigWithPassword("password1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mongodb_db_user.test", "name", "testuser"),
				),
			},
			{
				Config: testAccMongoDBUserConfigWithPassword("password2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mongodb_db_user.test", "name", "testuser"),
				),
			},
		},
	})
}

func testAccMongoDBUserConfig(role string) string {
	return fmt.Sprintf(`
provider "mongodb" {
  host          = "%s"
  port          = "%s"
  username      = "%s"
  password      = "%s"
  auth_database = "admin"
}

resource "mongodb_db_user" "test" {
  auth_database = "admin"
  name          = "testuser"
  password      = "testpassword"
  role {
    role = "%s"
    db   = "admin"
  }
}
`, os.Getenv("MONGODB_HOST"), getPortOrDefault(), os.Getenv("MONGODB_USERNAME"), os.Getenv("MONGODB_PASSWORD"), role)
}

func testAccMongoDBUserConfigWithPassword(password string) string {
	return fmt.Sprintf(`
provider "mongodb" {
  host          = "%s"
  port          = "%s"
  username      = "%s"
  password      = "%s"
  auth_database = "admin"
}

resource "mongodb_db_user" "test" {
  auth_database = "admin"
  name          = "testuser"
  password      = "%s"
  role {
    role = "read"
    db   = "admin"
  }
}
`, os.Getenv("MONGODB_HOST"), getPortOrDefault(), os.Getenv("MONGODB_USERNAME"), os.Getenv("MONGODB_PASSWORD"), password)
}

func getPortOrDefault() string {
	if port := os.Getenv("MONGODB_PORT"); port != "" {
		return port
	}
	return "27017"
}
