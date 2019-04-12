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
			"linux_group": groupResource(),
			"linux_user":  userResource(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Host:       d.Get("host").(string),
		Port:       d.Get("port").(int),
		User:       d.Get("user").(string),
		Password:   d.Get("password").(string),
		PrivateKey: os.ExpandEnv(d.Get("private_key").(string)),
	}

	log.Println("Initializing SSH client")
	return config.Client()
}
