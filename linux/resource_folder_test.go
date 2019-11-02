package linux

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccFolderCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: folderCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_folder.testfolder", "path", "/etc/testfolder"),
					resource.TestCheckNoResourceAttr("linux_folder.testfolder", "content"),
				),
			},
		},
	})
}

func TestAccFolderWithOwnerCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: folderWithOwnerCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_folder.testfolder", "path", "/etc/testfolder"),
					resource.TestCheckNoResourceAttr("linux_folder.testfolder", "content"),
					resource.TestCheckResourceAttr("linux_folder.testfolder", "owner", "testuser:testuser"),
				),
			},
		},
	})
}

func TestAccFolderWithPermissionsCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: folderWithPermissionsCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_folder.testfolder", "path", "/etc/testfolder"),
					resource.TestCheckNoResourceAttr("linux_folder.testfolder", "content"),
					resource.TestCheckResourceAttr("linux_folder.testfolder", "permissions", "777"),
				),
			},
		},
	})
}

func TestAccFolderWithAllAttrsCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: folderWithAllAttrsCreationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_folder.testfolder", "path", "/etc/testfolder"),
					resource.TestCheckNoResourceAttr("linux_folder.testfolder", "content"),
					resource.TestCheckResourceAttr("linux_folder.testfolder", "owner", "testuser:testuser"),
					resource.TestCheckResourceAttr("linux_folder.testfolder", "permissions", "777"),
				),
			},
		},
	})
}

func TestAccFolderUpdation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: folderWithAllAttrsCreationConfig,
			},
			resource.TestStep{
				Config: folderWithOwnerUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_folder.testfolder", "path", "/etc/testfolder"),
					resource.TestCheckNoResourceAttr("linux_folder.testfolder", "content"),
					resource.TestCheckResourceAttr("linux_folder.testfolder", "owner", "testuser_alt:testuser_alt"),
					resource.TestCheckResourceAttr("linux_folder.testfolder", "permissions", "777"),
				),
			},
			resource.TestStep{
				Config: folderWithOwnerPermissionsUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_folder.testfolder", "path", "/etc/testfolder"),
					resource.TestCheckNoResourceAttr("linux_folder.testfolder", "content"),
					resource.TestCheckResourceAttr("linux_folder.testfolder", "owner", "testuser_alt:testuser_alt"),
					resource.TestCheckResourceAttr("linux_folder.testfolder", "permissions", "666"),
				),
			},
			resource.TestStep{
				Config: folderWithOwnerPermissionsPathUpdatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_folder.testfolder", "path", "/etc/testfolder2"),
					resource.TestCheckNoResourceAttr("linux_folder.testfolder", "content"),
					resource.TestCheckResourceAttr("linux_folder.testfolder", "owner", "testuser_alt:testuser_alt"),
					resource.TestCheckResourceAttr("linux_folder.testfolder", "permissions", "666"),
				),
			},
			resource.TestStep{
				Config: folderWithOwnerPermissionsPathUpdatedAltConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linux_folder.testfolder", "path", "/etc/testfolder3"),
					resource.TestCheckNoResourceAttr("linux_folder.testfolder", "content"),
					resource.TestCheckResourceAttr("linux_folder.testfolder", "owner", "testuser_alt_alt:testuser_alt_alt"),
					resource.TestCheckResourceAttr("linux_folder.testfolder", "permissions", "766"),
				),
			},
		},
	})
}

const folderCreationConfig = `
resource "linux_folder" "testfolder" {
  path = "/etc/testfolder"
}
`
const folderWithOwnerCreationConfig = `
resource "linux_user" "testuser" {
	name = "testuser"
	uid = 1024
}
resource "linux_folder" "testfolder" {
  path = "/etc/testfolder"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
}
`
const folderWithPermissionsCreationConfig = `
resource "linux_folder" "testfolder" {
  path = "/etc/testfolder"
  permissions = 777
}
`
const folderWithAllAttrsCreationConfig = `
resource "linux_user" "testuser" {
	name = "testuser"
	uid = 1024
}
resource "linux_folder" "testfolder" {
  path = "/etc/testfolder"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 777
}
`
const folderWithOwnerUpdatedConfig = `
resource "linux_user" "testuser" {
	name = "testuser_alt"
	uid = 1025
}
resource "linux_folder" "testfolder" {
  path = "/etc/testfolder"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 777
}
`
const folderWithOwnerPermissionsUpdatedConfig = `
resource "linux_user" "testuser" {
	name = "testuser_alt"
	uid = 1025
}
resource "linux_folder" "testfolder" {
  path = "/etc/testfolder"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 666
}
`
const folderWithOwnerPermissionsPathUpdatedConfig = `
resource "linux_user" "testuser" {
	name = "testuser_alt"
	uid = 1025
}
resource "linux_folder" "testfolder" {
  path = "/etc/testfolder2"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 666
}
`
const folderWithOwnerPermissionsPathUpdatedAltConfig = `
resource "linux_user" "testuser" {
	name = "testuser_alt_alt"
	uid = 1026
}
resource "linux_folder" "testfolder" {
  path = "/etc/testfolder3"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 766
}
`
