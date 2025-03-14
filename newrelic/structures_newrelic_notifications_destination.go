package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"
)

func expandNotificationDestination(d *schema.ResourceData) (*notifications.AiNotificationsDestinationInput, error) {
	destination := notifications.AiNotificationsDestinationInput{
		Name: d.Get("name").(string),
		Type: notifications.AiNotificationsDestinationType(d.Get("type").(string)),
	}

	if attr, ok := d.GetOk("auth_basic"); ok {
		destination.Auth = expandNotificationDestinationAuthBasic(attr.([]interface{}))
	}

	if attr, ok := d.GetOk("auth_token"); ok {
		destination.Auth = expandNotificationDestinationAuthToken(attr.([]interface{}))
	}

	properties := d.Get("property")
	props := properties.(*schema.Set).List()
	destination.Properties = expandNotificationDestinationProperties(props)

	return &destination, nil
}

func expandNotificationDestinationAuthBasic(authRaw []interface{}) *notifications.AiNotificationsCredentialsInput {
	authInput := notifications.AiNotificationsCredentialsInput{}
	authInput.Type = notifications.AiNotificationsAuthTypeTypes.BASIC

	for _, a := range authRaw {
		aa := a.(map[string]interface{})
		authInput.Basic.User = aa["user"].(string)
		authInput.Basic.Password = notifications.SecureValue(aa["password"].(string))
	}

	return &authInput
}

func expandNotificationDestinationAuthToken(authRaw []interface{}) *notifications.AiNotificationsCredentialsInput {
	authInput := notifications.AiNotificationsCredentialsInput{}
	authInput.Type = notifications.AiNotificationsAuthTypeTypes.TOKEN

	for _, a := range authRaw {
		aa := a.(map[string]interface{})
		authInput.Token.Token = notifications.SecureValue(aa["token"].(string))
		authInput.Token.Prefix = aa["prefix"].(string)
	}

	return &authInput
}

func expandNotificationDestinationUpdate(d *schema.ResourceData) (*notifications.AiNotificationsDestinationUpdate, error) {
	destination := notifications.AiNotificationsDestinationUpdate{
		Name:   d.Get("name").(string),
		Active: d.Get("active").(bool),
	}

	if attr, ok := d.GetOk("auth_basic"); ok {
		destination.Auth = expandNotificationDestinationAuthBasic(attr.([]interface{}))
	}

	if attr, ok := d.GetOk("auth_token"); ok {
		destination.Auth = expandNotificationDestinationAuthToken(attr.([]interface{}))
	}

	properties := d.Get("property")
	props := properties.(*schema.Set).List()
	destination.Properties = expandNotificationDestinationProperties(props)

	return &destination, nil
}

func expandNotificationDestinationProperties(properties []interface{}) []notifications.AiNotificationsPropertyInput {
	props := []notifications.AiNotificationsPropertyInput{}

	for _, p := range properties {
		props = append(props, expandNotificationDestinationProperty(p.(map[string]interface{})))
	}

	return props
}

func expandNotificationDestinationProperty(cfg map[string]interface{}) notifications.AiNotificationsPropertyInput {
	property := notifications.AiNotificationsPropertyInput{}

	if propertyKey, ok := cfg["key"]; ok {
		property.Key = propertyKey.(string)
	}

	if propertyValue, ok := cfg["value"]; ok {
		property.Value = propertyValue.(string)
	}

	if propertyDisplayValue, ok := cfg["display_value"]; ok {
		property.DisplayValue = propertyDisplayValue.(string)
	}

	if propertyLabel, ok := cfg["label"]; ok {
		property.Label = propertyLabel.(string)
	}

	return property
}

func flattenNotificationDestination(destination *notifications.AiNotificationsDestination, d *schema.ResourceData) error {
	if destination == nil {
		return nil
	}

	var err error

	if err = d.Set("name", destination.Name); err != nil {
		return err
	}

	if err = d.Set("type", destination.Type); err != nil {
		return err
	}

	auth := flattenNotificationDestinationAuth(destination.Auth, d)

	authAttr := "auth_basic"
	switch destination.Auth.AuthType {
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.OAUTH2):
		authAttr = "auth_oauth2"
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.TOKEN):
		authAttr = "auth_token"
	}

	if err := d.Set(authAttr, auth); err != nil {
		return fmt.Errorf("[DEBUG] Error setting notification auth: %#v", err)
	}

	if err := d.Set("property", flattenNotificationDestinationProperties(destination.Properties)); err != nil {
		return err
	}

	if err := d.Set("active", destination.Active); err != nil {
		return err
	}

	if err := d.Set("account_id", destination.AccountID); err != nil {
		return err
	}

	if err := d.Set("is_user_authenticated", destination.IsUserAuthenticated); err != nil {
		return err
	}

	if err := d.Set("status", destination.Status); err != nil {
		return err
	}

	if err := d.Set("last_sent", destination.LastSent); err != nil {
		return err
	}

	return nil
}

func flattenNotificationDestinationAuth(a ai.AiNotificationsAuth, d *schema.ResourceData) []map[string]interface{} {
	authConfig := make([]map[string]interface{}, 1)

	switch a.AuthType {
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.BASIC):
		authConfig[0] = map[string]interface{}{
			"user":     a.User,
			"password": d.Get("auth_basic.0.password"),
		}
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.TOKEN):
		authConfig[0] = map[string]interface{}{
			"prefix": a.Prefix,
			"token":  d.Get("auth_token.0.token"),
		}
	case ai.AiNotificationsAuthType(notifications.AiNotificationsAuthTypeTypes.OAUTH2):
		// ...
	}

	return authConfig
}

func flattenNotificationDestinationProperties(p []notifications.AiNotificationsProperty) []map[string]interface{} {
	properties := []map[string]interface{}{}

	for _, property := range p {
		properties = append(properties, flattenNotificationDestinationProperty(property))
	}

	return properties
}

func flattenNotificationDestinationProperty(p notifications.AiNotificationsProperty) map[string]interface{} {
	propertyResult := make(map[string]interface{})

	propertyResult["key"] = p.Key
	propertyResult["value"] = p.Value
	propertyResult["display_value"] = p.DisplayValue
	propertyResult["label"] = p.Label

	return propertyResult
}
