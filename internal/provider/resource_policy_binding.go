package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePolicyBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyBindingCreate,
		ReadContext:   resourcePolicyBindingRead,
		UpdateContext: resourcePolicyBindingUpdate,
		DeleteContext: resourcePolicyBindingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the object this binding should apply to",
			},

			"policy": {
				Type:        schema.TypeString,
				Description: "UUID of the policy",
				Optional:    true,
			},
			"user": {
				Type:        schema.TypeString,
				Description: "UID of the user",
				Optional:    true,
			},
			"group": {
				Type:        schema.TypeString,
				Description: "UUID of the group",
				Optional:    true,
			},

			// General attributes
			"order": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"negate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
		},
	}
}

func resourcePolicyBindingSchemaToModel(d *schema.ResourceData) (*api.PolicyBindingRequest, diag.Diagnostics) {
	m := api.PolicyBindingRequest{
		Target:  d.Get("target").(string),
		Order:   int32(d.Get("order").(int)),
		Negate:  boolToPointer(d.Get("negate").(bool)),
		Enabled: boolToPointer(d.Get("enabled").(bool)),
		Timeout: intToPointer(d.Get("timeout").(int)),
	}

	if u, uSet := d.GetOk("policy"); uSet {
		m.Policy.Set(stringToPointer(u.(string)))
	} else {
		m.Policy.Set(nil)
	}

	if u, uSet := d.GetOk("user"); uSet {
		m.User.Set(intToPointer(u.(int)))
	} else {
		m.User.Set(nil)
	}

	if u, uSet := d.GetOk("group"); uSet {
		m.Group.Set(stringToPointer(u.(string)))
	} else {
		m.Group.Set(nil)
	}

	return &m, nil
}

func resourcePolicyBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourcePolicyBindingSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.PoliciesApi.PoliciesBindingsCreate(ctx).PolicyBindingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyBindingRead(ctx, d, m)
}

func resourcePolicyBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesBindingsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("target", res.Target)
	d.Set("policy", res.Policy)
	d.Set("user", res.User)
	d.Set("group", res.Group)
	d.Set("order", res.Order)
	d.Set("negate", res.Negate)
	d.Set("enabled", res.Enabled)
	d.Set("timeout", res.Timeout)
	return diags
}

func resourcePolicyBindingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourcePolicyBindingSchemaToModel(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.PoliciesApi.PoliciesBindingsUpdate(ctx, d.Id()).PolicyBindingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyBindingRead(ctx, d, m)
}

func resourcePolicyBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesBindingsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
