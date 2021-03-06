package sops

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"go.mozilla.org/sops/decrypt"
	"gopkg.in/yaml.v2"
)

// readData consolidates the logic of extracting the from the various input methods and setting it on the ResourceData
func readData(content []byte, format string, d *schema.ResourceData) error {
	cleartext, err := decrypt.Data(content, format)
	if err != nil {
		return err
	}

	if format != "raw" {
		var data map[string]interface{}
		switch format {
		case "json":
			err = json.Unmarshal(cleartext, &data)
		case "yaml":
			err = yaml.Unmarshal(cleartext, &data)
		default:
			return fmt.Errorf("Don't know how to unmarshal format %s", format)
		}
		if err != nil {
			return err
		}

		err = d.Set("data", flatten(data))
		if err != nil {
			return err
		}
	} else {
		err = d.Set("raw", string(cleartext))
		if err != nil {
			return err
		}
	}

	d.SetId("-")
	return nil
}
