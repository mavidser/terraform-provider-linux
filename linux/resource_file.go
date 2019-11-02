package linux

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func fileResource() *schema.Resource {
	return &schema.Resource{
		Create: fileResourceCreateWrapper(false),
		Read:   fileResourceReadWrapper(false),
		Update: fileResourceUpdateWrapper(false),
		Delete: fileResourceDelete,

		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePath,
			},
			"owner": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateOwner,
			},
			"permissions": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func createFile(client *Client, path string, isFolder bool) error {
	var command string
	if isFolder {
		command = "mkdir -p"
	} else {
		command = "touch"
	}
	command = fmt.Sprintf("%s %s", command, path)
	_, _, err := runCommand(client, false, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func applyOwner(client *Client, path string, owner string) error {
	command := fmt.Sprintf("chown %s %s", owner, path)
	_, _, err := runCommand(client, false, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func applyPermissions(client *Client, path string, permissions int) error {
	command := fmt.Sprintf("chmod %d %s", permissions, path)
	_, _, err := runCommand(client, false, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func writeContent(client *Client, path string, content string) error {
	command := fmt.Sprintf("cat > %s", path)
	_, _, err := runCommand(client, false, command, content)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func rollback(client *Client, err error, errMsg string, path string) error {
	err2 := errors.Wrap(err, errMsg)
	if err3 := deleteFile(client, path); err3 != nil {
		err3 = errors.Wrap(err2, err3.Error())
		return errors.Wrap(err3, "Couldn't delete file.")
	}
	return err2
}

func parsePermissionString(perms string) int {
	permMap := map[string]int{
		"---": 0,
		"--x": 1,
		"-w-": 2,
		"-wx": 3,
		"r--": 4,
		"r-x": 5,
		"rw-": 6,
		"rwx": 7,
	}
	return (permMap[perms[1:4]] * 100) +
		(permMap[perms[4:7]] * 10) +
		(permMap[perms[7:10]] * 1)
}

func fileResourceCreateWrapper(isFolder bool) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		client := m.(*Client)
		path := d.Get("path").(string)
		owner := d.Get("owner").(string)
		permissions := d.Get("permissions").(int)

		if err := createFile(client, path, isFolder); err != nil {
			return errors.Wrap(err, "Couldn't create file")
		}

		if owner != "" {
			if err := applyOwner(client, path, owner); err != nil {
				return rollback(client, err, "Couldn't apply owner, rolling back file creation", path)
			}
		}

		if permissions != 0 {
			if err := applyPermissions(client, path, permissions); err != nil {
				return rollback(client, err, "Couldn't apply permissions, rolling back file creation", path)
			}
		}

		if !isFolder {
			content := d.Get("content").(string)
			if content != "" {
				if err := writeContent(client, path, content); err != nil {
					return rollback(client, err, "Couldn't write content, rolling back file creation", path)
				}
			}
		}

		d.SetId(path)
		return fileResourceReadWrapper(isFolder)(d, m)
	}
}

func getDetails(client *Client, path string) (string, int, error) {
	command := fmt.Sprintf("ls -ld %s", path)
	stdout, _, err := runCommand(client, false, command, "")
	if err != nil {
		return "", 0, errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	if stdout == "" {
		return "", 0, fmt.Errorf("File not found with path %v", path)
	}
	fields := strings.Fields(stdout)
	permissions, user, group := parsePermissionString(fields[0]), fields[2], fields[3]
	if err != nil {
		return "", 0, errors.Wrap(err, fmt.Sprintf("Unable to parse permission string from %s", command))
	}
	return fmt.Sprintf("%s:%s", user, group), permissions, nil
}

func readFile(client *Client, path string) (string, error) {
	command := fmt.Sprintf("cat %s", path)
	stdout, _, err := runCommand(client, false, command, "")
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return stdout, nil
}

func fileResourceReadWrapper(isFolder bool) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		client := m.(*Client)
		id := d.Id()

		owner, permissions, err := getDetails(client, id)
		if err != nil {
			if strings.Contains(err.Error(), "File not found with path") {
				d.SetId("")
				return nil
			}
			return errors.Wrap(err, "Unable to ls the file")
		}

		if !isFolder {
			content, err := readFile(client, id)
			if err != nil {
				return errors.Wrap(err, "Unable to read the file")
			}
			d.Set("content", content)
		}

		d.Set("owner", owner)
		d.Set("permissions", permissions)
		return nil
	}
}

func moveFile(client *Client, oldPath string, newPath string) error {
	command := fmt.Sprintf("mv %s %s", oldPath, newPath)
	_, _, err := runCommand(client, false, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func fileResourceUpdateWrapper(isFolder bool) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		client := m.(*Client)

		path := d.Get("path").(string)
		owner := d.Get("owner").(string)
		permissions := d.Get("permissions").(int)

		oldPath := d.Id()
		oldOwner, oldPermissions, err := getDetails(client, oldPath)
		if err != nil {
			return errors.Wrap(err, "Unable to ls the file")
		}

		if !isFolder {
			content := d.Get("content").(string)
			oldContent, err := readFile(client, oldPath)
			if err != nil {
				return errors.Wrap(err, "Unable to read the file")
			}
			if oldContent != content {
				if err := writeContent(client, oldPath, content); err != nil {
					return errors.Wrap(err, "Couldn't rewrite content")
				}
			}
		}

		if oldPath != path {
			if err := moveFile(client, oldPath, path); err != nil {
				return errors.Wrap(err, "Couldn't mv file")
			}
			d.SetId(path)
		}

		if oldOwner != owner {
			if err := applyOwner(client, path, owner); err != nil {
				return errors.Wrap(err, "Couldn't apply owner")
			}
		}

		if oldPermissions != permissions {
			if err := applyPermissions(client, path, permissions); err != nil {
				return errors.Wrap(err, "Couldn't apply permissions")
			}
		}

		return fileResourceReadWrapper(isFolder)(d, m)
	}
}

func deleteFile(client *Client, path string) error {
	command := fmt.Sprintf("rm -rf %s", path)
	_, _, err := runCommand(client, false, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func fileResourceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	id := d.Id()

	return deleteFile(client, id)
}
