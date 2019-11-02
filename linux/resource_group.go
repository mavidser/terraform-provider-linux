package linux

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func groupResource() *schema.Resource {
	return &schema.Resource{
		Create: groupResourceCreate,
		Read:   groupResourceRead,
		Update: groupResourceUpdate,
		Delete: groupResourceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gid": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"system": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func groupResourceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	name := d.Get("name").(string)
	gid := d.Get("gid").(int)
	system := d.Get("system").(bool)

	err := createGroup(client, name, gid, system)
	if err != nil {
		return errors.Wrap(err, "Couldn't create group")
	}

	gid, err = getGroupId(client, name)
	if err != nil {
		return errors.Wrap(err, "Couldn't get gid")
	}

	d.Set("gid", gid)

	d.SetId(fmt.Sprintf("%v", gid))
	return groupResourceRead(d, m)
}

func createGroup(client *Client, name string, gid int, system bool) error {
	command := "/usr/sbin/groupadd"

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

func getGroupId(client *Client, name string) (int, error) {
	command := fmt.Sprintf("getent group %s", name)
	stdout, _, err := runCommand(client, false, command, "")
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	if stdout == "" {
		return 0, fmt.Errorf("Group not found with name %v", name)
	}
	gid, err := strconv.Atoi(strings.Split(stdout, ":")[2])
	if err != nil {
		return 0, err
	}
	return gid, nil
}

func getGroupName(client *Client, gid int) (string, error) {
	command := fmt.Sprintf("getent group %d", gid)
	stdout, _, err := runCommand(client, false, command, "")
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	if stdout == "" {
		return "", fmt.Errorf("Group not found with id %v", gid)
	}
	name := strings.Split(stdout, ":")[0]
	return name, nil
}

func groupResourceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	gid, err := strconv.Atoi(d.Id())
	if err != nil {
		return errors.Wrap(err, "ID stored is not int")
	}
	name, err := getGroupName(client, gid)
	if err != nil {
		log.Printf("%v", err)
		log.Printf("Error getting group name, will recreate it")
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	return nil
}

func groupResourceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	gid, err := strconv.Atoi(d.Id())
	if err != nil {
		return errors.Wrap(err, "ID stored is not int")
	}
	name := d.Get("name").(string)
	oldname, err := getGroupName(client, gid)
	if err != nil {
		return errors.Wrap(err, "Failed to get group name")
	}

	if oldname != name {
		command := fmt.Sprintf("/usr/sbin/groupmod %s -n %s", oldname, name)
		_, _, err = runCommand(client, true, command, "")
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
		}
	}
	return groupResourceRead(d, m)
}

func deleteGroup(client *Client, name string) error {
	command := fmt.Sprintf("/usr/sbin/groupdel %s", name)
	_, _, err := runCommand(client, true, command, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Command failed: %s", command))
	}
	return nil
}

func groupResourceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	gid, err := strconv.Atoi(d.Id())
	if err != nil {
		return errors.Wrap(err, "ID stored is not int")
	}
	name, err := getGroupName(client, gid)
	if err != nil {
		return errors.Wrap(err, "Failed to get group name")
	}

	return deleteGroup(client, name)
}
