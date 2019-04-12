package linux

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func groupResource() *schema.Resource {
	return &schema.Resource{
		Create: groupResourceServerCreate,
		Read:   groupResourceServerRead,
		Update: groupResourceServerUpdate,
		Delete: groupResourceServerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gid": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"system": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func groupResourceServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	name := d.Get("name").(string)
	gid := d.Get("gid").(int)
	system := d.Get("system").(bool)

	err := createGroup(client, name, gid, system)
	if err != nil {
		return err
	}

	gid, err = getGroupId(client, name)
	if err != nil {
		return err
	}

	d.Set("gid", gid)

	d.SetId(name)
	return groupResourceServerRead(d, m)
}

func createGroup(client *Client, name string, gid int, system bool) error {
	command := "sudo /usr/sbin/groupadd"

	if gid > 0 {
		command = fmt.Sprintf("%s --gid %d", command, gid)
	}
	if system {
		command = fmt.Sprintf("%s --system", command)
	}
	command = fmt.Sprintf("%s %s", command, name)

	session, err := client.connection.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create session: %s", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stderr for session: %v", err)
	}

	log.Printf("Running command %s", command)

	err = session.Run(command)
	if err != nil {
		output, err2 := ioutil.ReadAll(stderr)
		if err2 != nil {
			log.Printf("Unable to read stderr for command: %v", err)
		}
		log.Printf("Stderr output: %s", string(output))

		return fmt.Errorf("Error running command %s: %s", command, err)
	}
	return nil
}

func getGroupId(client *Client, name string) (int, error) {
	command := fmt.Sprintf("getent group %s", name)
	session, err := client.connection.NewSession()
	if err != nil {
		return 0, fmt.Errorf("Failed to create session: %s", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return 0, fmt.Errorf("Unable to setup stderr for session: %v", err)
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return 0, fmt.Errorf("Unable to setup stdout for session: %v", err)
	}

	log.Printf("Running command %s", command)

	err = session.Run(command)
	if err != nil {
		output, err2 := ioutil.ReadAll(stderr)
		if err2 != nil {
			log.Printf("Unable to read stderr for command: %v", err)
		}
		log.Printf("Stderr output: %s", strings.TrimSpace(string(output)))

		return 0, fmt.Errorf("Error running command %s: %s", command, err)
	}

	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		return 0, fmt.Errorf("Unable to read stdout for command: %v", err)
	}
	gid, err := strconv.Atoi(strings.Split(string(output), ":")[2])
	if err != nil {
		return 0, err
	}
	return gid, nil
}

func groupResourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func groupResourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func groupResourceServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	name := d.Id()

	command := fmt.Sprintf("sudo /usr/sbin/groupdel %s", name)
	session, err := client.connection.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create session: %s", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stderr for session: %v", err)
	}

	log.Printf("Running command %s", command)

	err = session.Run(command)
	if err != nil {
		output, err2 := ioutil.ReadAll(stderr)
		if err2 != nil {
			log.Printf("Unable to read stderr for command: %v", err)
		}
		log.Printf("Stderr output: %s", string(output))

		return fmt.Errorf("Error running command %s: %s", command, err)
	}

	return nil
}
