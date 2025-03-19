package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBlueprint(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlueprintSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_blueprint.default", "name", "Default - Brand"),
					resource.TestCheckResourceAttr("data.authentik_blueprint.default", "path", "default/default-brand.yaml"),
				),
			},
		},
	})
}

const testAccDataSourceBlueprintSimple = `
data "authentik_blueprint" "default" {
  name = "Default - Brand"
}
`
