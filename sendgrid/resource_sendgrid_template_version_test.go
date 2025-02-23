package sendgrid_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sendgrid "github.com/taharah/terraform-provider-sendgrid/sdk"
)

func TestAccSendgridTemplateVersionBasic(t *testing.T) {
	templateName := "terraform-template-" + acctest.RandString(10)
	templateVersionName := "terraform-template-version-" + acctest.RandString(10)
	subject := "terraform-subject-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridTemplateVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridTemplateVersionConfigBasic(
					templateName,
					templateVersionName,
					subject,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridTemplateVersionExists("sendgrid_template_version.this"),
				),
			},
		},
	})
}

func testAccCheckSendgridTemplateVersionDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*sendgrid.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sendgrid_template_version" {
			continue
		}

		templateID := rs.Primary.Attributes["template_id"]
		id := rs.Primary.ID

		if _, err := c.DeleteTemplateVersion(templateID, id); err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckSendgridTemplateVersionConfigBasic(
	templateName, templateVersionName, subject string,
) string {
	return fmt.Sprintf(`
resource "sendgrid_template" "this" {
  name = %q
}
resource "sendgrid_template_version" "this" {
  template_id = sendgrid_template.this.id
  name = %q
  subject = %q
}`, templateName, templateVersionName, subject)
}

func testAccCheckSendgridTemplateVersionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No templateVersionID set")
		}

		return nil
	}
}
