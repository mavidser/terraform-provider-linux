package linux

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccGroupCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: groupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_group.testgroup", "name", "testgroup"),
					testAccCheckGID("testgroup", func(gid int) error { return nil }),
				),
			},
		},
	})
}

func TestAccSystemGroupCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: systemGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_group.testgroup", "name", "testgroup"),
					resource.TestCheckResourceAttr("linux_group.testgroup", "system", "true"),
					testAccCheckGID("testgroup", func(gid int) error {
						if gid > 1000 {
							return fmt.Errorf("System group gid should be less than 1000")
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestAccGroupWithGIDCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: groupWithGIDConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_group.testgroup", "name", "testgroup"),
					resource.TestCheckResourceAttr("linux_group.testgroup", "gid", "1024"),
					testAccCheckGID("testgroup", func(gid int) error {
						if gid != 1024 {
							return fmt.Errorf("GID should be 1024")
						}
						return nil
					}),
				),
			},
		},
	})
}

func testAccCheckGID(groupname string, check func(int) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		gid, err := getGroupId(client, groupname)
		if err != nil {
			return err
		}
		return check(gid)
	}
}

func TestAccGroupUpdation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: groupWithGIDConfig,
			},
			resource.TestStep{
				Config: groupWithNameUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_group.testgroup", "name", "testgroup_alt"),
					resource.TestCheckResourceAttr("linux_group.testgroup", "gid", "1024"),
					testAccCheckGID("testgroup_alt", func(gid int) error {
						if gid != 1024 {
							return fmt.Errorf("GID should be 1024")
						}
						return nil
					}),
				),
			},
			resource.TestStep{
				Config: groupWithNameGIDUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_group.testgroup", "name", "testgroup_alt"),
					resource.TestCheckResourceAttr("linux_group.testgroup", "gid", "1025"),
				),
			},
		},
	})
}

const groupConfig = `
resource "linux_group" "testgroup" {
	name = "testgroup"
}
`
const systemGroupConfig = `
resource "linux_group" "testgroup" {
	name = "testgroup"
	system = true
}
`
const groupWithGIDConfig = `
resource "linux_group" "testgroup" {
	name = "testgroup"
	gid = 1024
}
`
const groupWithNameUpdatedConfig = `
resource "linux_group" "testgroup" {
	name = "testgroup_alt"
	gid = 1024
}
`
const groupWithNameGIDUpdatedConfig = `
resource "linux_group" "testgroup" {
	name = "testgroup_alt"
	gid = 1025
}
`
