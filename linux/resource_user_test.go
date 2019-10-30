package linux

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	// "github.com/pkg/errors"
)

func TestAccUserCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_user.testuser", "name", "testuser"),
					testAccCheckUID("testuser", func(uid int) error { return nil }),
				),
			},
		},
	})
}

func TestAccSystemUserCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSystemUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_user.testuser", "name", "testuser"),
					resource.TestCheckResourceAttr("linux_user.testuser", "system", "true"),
					testAccCheckUID("testuser", func(uid int) error {
						if uid > 1000 {
							return fmt.Errorf("System user uid should be less than 1000")
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestAccUserWithUIDCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccUserWithUIDConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_user.testuser", "name", "testuser"),
					resource.TestCheckResourceAttr("linux_user.testuser", "uid", "1024"),
					testAccCheckUID("testuser", func(uid int) error {
						if uid != 1024 {
							return fmt.Errorf("UID should be 1024")
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestAccUserWithGroupCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccUserWithGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_user.testuser", "name", "testuser"),
					resource.TestCheckResourceAttr("linux_user.testuser", "uid", "1024"),
					resource.TestCheckResourceAttr("linux_user.testuser", "uid", "1024"),
					resource.TestCheckResourceAttr("linux_user.testuser", "gid", "1048"),
					testAccCheckUID("testuser", func(uid int) error {
						if uid != 1024 {
							return fmt.Errorf("UID should be 1024")
						}
						return nil
					}),
					testAccCheckGID("testgroup", func(gid int) error {
						if gid != 1048 {
							return fmt.Errorf("GID should be 1048")
						}
						return nil
					}),
					testAccCheckGIDForUser("testuser", func(gid int) error {
						if gid != 1048 {
							return fmt.Errorf("GID should be 1048")
						}
						return nil
					}),
				),
			},
		},
	})
}

func testAccCheckUID(username string, check func(int) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		uid, err := getUserId(client, username)
		if err != nil {
			return err
		}
		return check(uid)
	}
}

func testAccCheckGIDForUser(username string, check func(int) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		gid, err := getGroupIdForUser(client, username)
		if err != nil {
			return err
		}
		return check(gid)
	}
}

const testAccUserConfig = `
resource "linux_user" "testuser" {
	name = "testuser"
}
`
const testAccSystemUserConfig = `
resource "linux_user" "testuser" {
	name = "testuser"
	system = true
}
`
const testAccUserWithUIDConfig = `
resource "linux_user" "testuser" {
	name = "testuser"
	uid = 1024
}
`
const testAccUserWithGroupConfig = `
resource "linux_group" "testgroup" {
	name = "testgroup"
	gid = 1048
}
resource "linux_user" "testuser" {
	name = "testuser"
	uid = 1024
	gid = linux_group.testgroup.gid
}
`
