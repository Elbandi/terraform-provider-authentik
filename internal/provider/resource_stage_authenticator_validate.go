package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api"
)

func resourceStageAuthenticatorValidate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStageAuthenticatorValidateCreate,
		ReadContext:   resourceStageAuthenticatorValidateRead,
		UpdateContext: resourceStageAuthenticatorValidateUpdate,
		DeleteContext: resourceStageAuthenticatorValidateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"not_configured_action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"device_classes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"configuration_stage": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceStageAuthenticatorValidateSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorValidateStageRequest {
	r := api.AuthenticatorValidateStageRequest{
		Name: d.Get("name").(string),
	}

	if h, hSet := d.GetOk("not_configured_action"); hSet {
		action := api.NotConfiguredActionEnum(h.(string))
		r.NotConfiguredAction = &action
	}
	if h, hSet := d.GetOk("configuration_stage"); hSet {
		r.ConfigurationStage.Set(stringToPointer(h.(string)))
	}

	classes := make([]api.DeviceClassesEnum, 0)
	for _, classesS := range d.Get("device_classes").([]interface{}) {
		classes = append(classes, api.DeviceClassesEnum(classesS.(string)))
	}
	r.DeviceClasses = &classes

	return &r
}

func resourceStageAuthenticatorValidateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorValidateSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorValidateCreate(ctx).AuthenticatorValidateStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorValidateRead(ctx, d, m)
}

func resourceStageAuthenticatorValidateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorValidateRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.Set("name", res.Name)
	d.Set("not_configured_action", res.NotConfiguredAction)
	if res.ConfigurationStage.IsSet() {
		d.Set("configuration_stage", res.ConfigurationStage.Get())
	}
	d.Set("device_classes", res.DeviceClasses)
	return diags
}

func resourceStageAuthenticatorValidateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorValidateSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorValidateUpdate(ctx, d.Id()).AuthenticatorValidateStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorValidateRead(ctx, d, m)
}

func resourceStageAuthenticatorValidateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorValidateDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
