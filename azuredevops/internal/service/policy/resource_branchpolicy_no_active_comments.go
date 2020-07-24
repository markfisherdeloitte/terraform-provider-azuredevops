package policy

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
)

type noActiveCommentsPolicySettings struct{}

func ResourceBranchPolicyNoActiveComments() *schema.Resource {
	return genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: noActiveCommentsFlattenFunc,
		ExpandFunc:  noActiveCommentsExpandFunc,
		PolicyType:  NoActiveComments,
	})
}

func noActiveCommentsFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}
	policyAsJSON, err := json.Marshal(policyConfig.Settings)
	if err != nil {
		return fmt.Errorf("Unable to marshal policy settings into JSON: %+v", err)
	}

	policySettings := autoReviewerPolicySettings{}
	err = json.Unmarshal(policyAsJSON, &policySettings)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal branch policy settings (%+v): %+v", policySettings, err)
	}

	settingsList := d.Get(SchemaSettings).([]interface{})
	settings := settingsList[0].(map[string]interface{})

	settings["allowNoFastForward"] = true
	settings["allowRebase"] = true
	settings["allowRebaseMerge"] = true
	settings["allowSquash"] = true

	d.Set(SchemaSettings, settingsList)
	return nil
}

func noActiveCommentsExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	return policyConfig, projectID, nil
}
