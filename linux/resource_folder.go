package linux

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func folderResource() *schema.Resource {
	return &schema.Resource{
		Create: fileResourceCreateWrapper(true),
		Read:   fileResourceReadWrapper(true),
		Update: fileResourceUpdateWrapper(true),
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
		},
	}
}
