package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePolicyExpression() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyExpressionRead,
		Description: "Customization --- Get Policy Expression by name",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execution_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourcePolicyExpressionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.PoliciesApi.PoliciesExpressionList(ctx)
	if name, ok := d.GetOk("name"); ok {
		req = req.Name(name.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching policy expression found")
	}
	f := res.Results[0]
	d.SetId(f.Pk)
	setWrapper(d, "name", f.Name)
	setWrapper(d, "execution_logging", f.ExecutionLogging)
	setWrapper(d, "expression", f.Expression)
	//	setWrapper(d, "uuid", f.Pk)
	return diags
}
