package linux

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
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
				ForceNew: true,
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

func userResourceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	name := d.Get("name").(string)
	uid := d.Get("uid").(int)
	gid := d.Get("gid").(int)
	system := d.Get("system").(bool)

	err := createUser(client, name, uid, gid, system)
	if err != nil {
		return err
	}

	uid, err = getUserId(client, name)
	if err != nil {
		return err
	}

	d.Set("uid", uid)

	d.SetId(name)
	return userResourceRead(d, m)
}

func createUser(client *Client, name string, uid int, gid int, system bool) error {
	command := "sudo /usr/sbin/useradd"

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

func getUserId(client *Client, name string) (int, error) {
	command := fmt.Sprintf("id --user %s", name)
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
		log.Printf("Stderr output: %s", string(output))

		return 0, fmt.Errorf("Error running command %s: %s", command, err)
	}

	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		return 0, fmt.Errorf("Unable to read stdout for command: %v", err)
	}
	uid, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func userResourceRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func userResourceUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func userResourceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	name := d.Id()

	command := fmt.Sprintf("sudo /usr/sbin/userdel %s", name)
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
