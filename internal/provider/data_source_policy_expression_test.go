package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePolicyExpression(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicyExpressionSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_policy_expression.inbuilt", "name", "default-authentication-flow-password-stage"),
					resource.TestCheckResourceAttrSet("data.authentik_policy_expression.inbuilt", "expression"),
				),
			},
		},
	})
}

const testAccDataSourcePolicyExpressionSimple = `
data "authentik_policy_expression" "inbuilt" {
  name = "default-authentication-flow-password-stage"
}
`
