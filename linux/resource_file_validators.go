package linux

import (
	"fmt"
	"strings"
)

func validatePath(vi interface{}, k string) (ws []string, errors []error) {
	v, err := vi.(string)
	if !err {
		errors = append(errors, fmt.Errorf("path should be a string"))
	}
	if !strings.HasPrefix(v, "/") {
		errors = append(errors, fmt.Errorf("path should be an absolute path"))
	}
	return
}

func validateOwner(vi interface{}, k string) (ws []string, errors []error) {
	v, err := vi.(string)
	if !err {
		errors = append(errors, fmt.Errorf("owner should be a string"))
	}
	if i := strings.Index(v, ":"); i <= 0 || i >= len(v)-1 {
		errors = append(errors, fmt.Errorf("owner should be of the form user:group"))
	}
	return
}
