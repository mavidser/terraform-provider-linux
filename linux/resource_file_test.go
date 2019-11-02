package linux

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccFileCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fileCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_file.testfile", "path", "/etc/testfile"),
					resource.TestCheckResourceAttr("linux_file.testfile", "content", ""),
				),
			},
		},
	})
}

func TestAccFileWithContentCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fileWithContentCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_file.testfile", "path", "/etc/testfile"),
					resource.TestCheckResourceAttr("linux_file.testfile", "content", "testcontent"),
				),
			},
		},
	})
}

func TestAccFileWithOwnerCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fileWithOwnerCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_file.testfile", "path", "/etc/testfile"),
					resource.TestCheckResourceAttr("linux_file.testfile", "content", ""),
					resource.TestCheckResourceAttr("linux_file.testfile", "owner", "testuser:testuser"),
				),
			},
		},
	})
}

func TestAccFileWithPermissionsCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fileWithPermissionsCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_file.testfile", "path", "/etc/testfile"),
					resource.TestCheckResourceAttr("linux_file.testfile", "content", ""),
					resource.TestCheckResourceAttr("linux_file.testfile", "permissions", "777"),
				),
			},
		},
	})
}

func TestAccFileWithAllAttrsCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fileWithAllAttrsCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_file.testfile", "path", "/etc/testfile"),
					resource.TestCheckResourceAttr("linux_file.testfile", "content", "testcontent"),
					resource.TestCheckResourceAttr("linux_file.testfile", "owner", "testuser:testuser"),
					resource.TestCheckResourceAttr("linux_file.testfile", "permissions", "777"),
				),
			},
		},
	})
}

func TestAccFileUpdation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fileWithAllAttrsCreationConfig,
			},
			resource.TestStep{
				Config: fileWithOwnerUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_file.testfile", "path", "/etc/testfile"),
					resource.TestCheckResourceAttr("linux_file.testfile", "content", "testcontent"),
					resource.TestCheckResourceAttr("linux_file.testfile", "owner", "testuser_alt:testuser_alt"),
					resource.TestCheckResourceAttr("linux_file.testfile", "permissions", "777"),
				),
			},
			resource.TestStep{
				Config: fileWithOwnerPermissionsUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_file.testfile", "path", "/etc/testfile"),
					resource.TestCheckResourceAttr("linux_file.testfile", "content", "testcontent"),
					resource.TestCheckResourceAttr("linux_file.testfile", "owner", "testuser_alt:testuser_alt"),
					resource.TestCheckResourceAttr("linux_file.testfile", "permissions", "666"),
				),
			},
			resource.TestStep{
				Config: fileWithOwnerPermissionsContentUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_file.testfile", "path", "/etc/testfile"),
					resource.TestCheckResourceAttr("linux_file.testfile", "content", "testcontent2"),
					resource.TestCheckResourceAttr("linux_file.testfile", "owner", "testuser_alt:testuser_alt"),
					resource.TestCheckResourceAttr("linux_file.testfile", "permissions", "666"),
				),
			},
			resource.TestStep{
				Config: fileWithOwnerPermissionsContentPathUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_file.testfile", "path", "/etc/testfile2"),
					resource.TestCheckResourceAttr("linux_file.testfile", "content", "testcontent2"),
					resource.TestCheckResourceAttr("linux_file.testfile", "owner", "testuser_alt:testuser_alt"),
					resource.TestCheckResourceAttr("linux_file.testfile", "permissions", "666"),
				),
			},
			resource.TestStep{
				Config: fileWithOwnerPermissionsContentPathUpdatedAltConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_file.testfile", "path", "/etc/testfile3"),
					resource.TestCheckResourceAttr("linux_file.testfile", "content", "testcontent3"),
					resource.TestCheckResourceAttr("linux_file.testfile", "owner", "testuser_alt_alt:testuser_alt_alt"),
					resource.TestCheckResourceAttr("linux_file.testfile", "permissions", "766"),
				),
			},
		},
	})
}

const fileCreationConfig = `
resource "linux_file" "testfile" {
  path = "/etc/testfile"
}
`
const fileWithContentCreationConfig = `
resource "linux_file" "testfile" {
  path = "/etc/testfile"
  content = "testcontent"
}
`
const fileWithOwnerCreationConfig = `
resource "linux_user" "testuser" {
	name = "testuser"
	uid = 1024
}
resource "linux_file" "testfile" {
  path = "/etc/testfile"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
}
`
const fileWithPermissionsCreationConfig = `
resource "linux_file" "testfile" {
  path = "/etc/testfile"
  permissions = 777
}
`
const fileWithAllAttrsCreationConfig = `
resource "linux_user" "testuser" {
	name = "testuser"
	uid = 1024
}
resource "linux_file" "testfile" {
  path = "/etc/testfile"
  content = "testcontent"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 777
}
`
const fileWithOwnerUpdatedConfig = `
resource "linux_user" "testuser" {
	name = "testuser_alt"
	uid = 1025
}
resource "linux_file" "testfile" {
  path = "/etc/testfile"
  content = "testcontent"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 777
}
`
const fileWithOwnerPermissionsUpdatedConfig = `
resource "linux_user" "testuser" {
	name = "testuser_alt"
	uid = 1025
}
resource "linux_file" "testfile" {
  path = "/etc/testfile"
  content = "testcontent"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 666
}
`
const fileWithOwnerPermissionsContentUpdatedConfig = `
resource "linux_user" "testuser" {
	name = "testuser_alt"
	uid = 1025
}
resource "linux_file" "testfile" {
  path = "/etc/testfile"
  content = "testcontent2"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 666
}
`
const fileWithOwnerPermissionsContentPathUpdatedConfig = `
resource "linux_user" "testuser" {
	name = "testuser_alt"
	uid = 1025
}
resource "linux_file" "testfile" {
  path = "/etc/testfile2"
  content = "testcontent2"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 666
}
`
const fileWithOwnerPermissionsContentPathUpdatedAltConfig = `
resource "linux_user" "testuser" {
	name = "testuser_alt_alt"
	uid = 1026
}
resource "linux_file" "testfile" {
  path = "/etc/testfile3"
  content = "testcontent3"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 766
}
`
