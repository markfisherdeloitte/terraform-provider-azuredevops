package policy

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
)

const (
	allowNoFastForward = "allow_no_fast_forward"
	allowRebase        = "allow_rebase"
	allowRebaseMerge   = "allow_rebase_merge"
	allowSquash        = "allow_squash"
)

type allowedMergeTypesPolicySettings struct {
	AllowNoFastForward bool `json:"allowNoFastForward"`
	AllowRebase        bool `json:"allowRebase"`
	AllowRebaseMerge   bool `json:"allowRebaseMerge"`
	AllowSquash        bool `json:"allowSquash"`
}

func ResourceBranchPolicyAllowedMergeTypes() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: allowedMergeTypesFlattenFunc,
		ExpandFunc:  allowedMergeTypesExpandFunc,
		PolicyType:  AllowedMergeTypes,
	})

	settingsSchema := resource.Schema[SchemaSettings].Elem.(*schema.Resource).Schema
	settingsSchema[allowNoFastForward] = &schema.Schema{
		Type:     schema.TypeBool,
		Default:  true,
		Optional: true,
	}
	settingsSchema[allowRebase] = &schema.Schema{
		Type:     schema.TypeBool,
		Default:  true,
		Optional: true,
	}
	settingsSchema[allowRebaseMerge] = &schema.Schema{
		Type:     schema.TypeBool,
		Default:  true,
		Optional: true,
	}
	settingsSchema[allowSquash] = &schema.Schema{
		Type:     schema.TypeBool,
		Default:  true,
		Optional: true,
	}
	return resource
}

func allowedMergeTypesFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}
	policyAsJSON, err := json.Marshal(policyConfig.Settings)
	if err != nil {
		return fmt.Errorf("Unable to marshal policy settings into JSON: %+v", err)
	}

	policySettings := allowedMergeTypesPolicySettings{}
	err = json.Unmarshal(policyAsJSON, &policySettings)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal branch policy settings (%+v): %+v", policySettings, err)
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	settings[allowNoFastForward] = policySettings.AllowNoFastForward
	settings[allowRebase] = policySettings.AllowRebase
	settings[allowRebaseMerge] = policySettings.AllowRebaseMerge
	settings[allowSquash] = policySettings.AllowSquash

	d.Set(SchemaSettings, settingsList)
	return nil
}

func allowedMergeTypesExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["allowNoFastForward"] = settings[allowNoFastForward].(bool)
	policySettings["allowRebase"] = settings[allowRebase].(bool)
	policySettings["allowRebaseMerge"] = settings[allowRebaseMerge].(bool)
	policySettings["allowSquash"] = settings[allowSquash].(bool)

	return policyConfig, projectID, nil
}
