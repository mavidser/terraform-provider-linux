package linux

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func userResource() *schema.Resource {
	return &schema.Resource{
		Create: userResourceCreate,
		Read:   userResourceRead,
		Update: userResourceUpdate,
		Delete: userResourceDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"uid": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"gid": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"system": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func userResourceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	name := d.Get("name").(string)
	uid := d.Get("uid").(int)
	gid := d.Get("gid").(int)
	system := d.Get("system").(bool)

	err := createUser(client, name, uid, gid, system)
	if err != nil {
		return errors.Wrap(err, "Couldn't create user")
	}

	uid, err = getUserId(client, name)
	if err != nil {
		return errors.Wrap(err, "Couldn't get uid")
	}

	d.Set("uid", uid)

	d.SetId(fmt.Sprintf("%v", uid))
	return userResourceRead(d, m)
}

func createUser(client *Client, name string, uid int, gid int, system bool) error {
	command := "/usr/sbin/useradd"

	if uid > 0 {
		command = fmt.Sprintf("%s --uid %d", command, uid)
	}
	if gid > 0 {
		command = fmt.Sprintf("%s --gid %d", command, gid)
	}
	if system {
		command = fmt.Sprintf("%s --system", command)
	}
	command = fmt.Sprintf("%s %s", command, name)
	_, _, err := runCommand(client, true, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func getUserId(client *Client, name string) (int, error) {
	command := fmt.Sprintf("id --user %s", name)
	stdout, _, err := runCommand(client, false, command, "")
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	if stdout == "" {
		return 0, fmt.Errorf("User not found with name %v", name)
	}
	uid, err := strconv.Atoi(strings.TrimSpace(stdout))
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func getUserName(client *Client, uid int) (string, error) {
	command := fmt.Sprintf("getent passwd %d", uid)
	stdout, _, err := runCommand(client, false, command, "")
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	if stdout == "" {
		return "", fmt.Errorf("User not found with id %v", uid)
	}
	name := strings.Split(stdout, ":")[0]
	return name, nil
}

func getGroupIdForUser(client *Client, name string) (int, error) {
	command := fmt.Sprintf("getent passwd %s", name)
	stdout, _, err := runCommand(client, false, command, "")
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	if stdout == "" {
		return 0, fmt.Errorf("Group not found for user %v", name)
	}
	uid, err := strconv.Atoi(strings.TrimSpace(strings.Split(stdout, ":")[3]))
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func userResourceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	uid, err := strconv.Atoi(d.Id())
	if err != nil {
		return errors.Wrap(err, "ID stored is not int")
	}
	name, err := getUserName(client, uid)
	if err != nil {
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	gid, err := getGroupIdForUser(client, name)
	if err != nil {
		return errors.Wrap(err, "Couldn't find group for user")
	}
	d.Set("gid", gid)
	return nil
}

func userResourceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	uid, err := strconv.Atoi(d.Id())
	if err != nil {
		return errors.Wrap(err, "ID stored is not int")
	}
	name := d.Get("name").(string)
	gid := d.Get("gid").(int)
	oldname, err := getUserName(client, uid)
	if err != nil {
		return errors.Wrap(err, "Failed to get user name")
	}
	oldgid, err := getGroupIdForUser(client, oldname)
	if err != nil {
		return errors.Wrap(err, "Failed to get user gid")
	}

	if oldname != name {
		command := fmt.Sprintf("/usr/sbin/usermod %s -l %s", oldname, name)
		_, _, err = runCommand(client, true, command, "")
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
		}
	}

	if oldgid != gid {
		command := fmt.Sprintf("/usr/sbin/usermod %s -g %d", name, gid)
		_, _, err = runCommand(client, true, command, "")
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
		}
	}
	return userResourceRead(d, m)
}

func userResourceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	uid, err := strconv.Atoi(d.Id())
	if err != nil {
		errors.Wrap(err, "ID stored is not int")
	}
	name, err := getUserName(client, uid)
	if err != nil {
		return errors.Wrap(err, "Failed to get user name")
	}

	command := fmt.Sprintf("/usr/sbin/userdel %s", name)
	_, _, err = runCommand(client, true, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}
