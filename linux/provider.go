package linux

import (
	"log"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_SSH_USER", ""),
				Description: "The Docker daemon address",
			},
			"use_sudo": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_USE_SUDO", ""),
				Description: "The Docker daemon address",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_SSH_HOST", ""),
				Description: "The Docker daemon address",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_SSH_PORT", 22),
				Description: "The Docker daemon address",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_SSH_PASSWORD", ""),
				Description: "The Docker daemon address",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_LINUX_SSH_PRIVATE_KEY", "$HOME/.ssh/id_rsa"),
				Description: "The Docker daemon address",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"linux_group":  groupResource(),
			"linux_user":   userResource(),
			"linux_file":   fileResource(),
			"linux_folder": folderResource(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	user := d.Get("user").(string)

	var useSudoBool bool
	useSudo, ok := d.GetOk("use_sudo")
	if !ok {
		if user == "root" {
			useSudoBool = false
		} else {
			useSudoBool = true
		}
	} else {
		useSudoBool = useSudo.(bool)
	}
	config := Config{
		Host:       d.Get("host").(string),
		Port:       d.Get("port").(int),
		User:       user,
		Password:   d.Get("password").(string),
		PrivateKey: os.ExpandEnv(d.Get("private_key").(string)),
		UseSudo:    useSudoBool,
	}

	log.Println("Initializing SSH client")
	return config.Client()
}
