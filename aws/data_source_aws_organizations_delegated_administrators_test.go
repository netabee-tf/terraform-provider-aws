package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDataSourceAwsOrganizationsDelegatedAdministrators_basic(t *testing.T) {
	var providers []*schema.Provider
	dataSourceName := "data.aws_organizations_delegated_administrators.test"
	servicePrincipal := "config-multiaccountsetup.amazonaws.com"
	dataSourceIdentity := "data.aws_caller_identity.delegated"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccAlternateAccountPreCheck(t)
		},
		ErrorCheck:        testAccErrorCheck(t, organizations.EndpointsID),
		ProviderFactories: testAccProviderFactoriesAlternate(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsOrganizationsDelegatedAdministratorsConfig(servicePrincipal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "delegated_administrators.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "delegated_administrators.0.id", dataSourceIdentity, "account_id"),
					testAccCheckResourceAttrRfc3339(dataSourceName, "delegated_administrators.0.delegation_enabled_date"),
					testAccCheckResourceAttrRfc3339(dataSourceName, "delegated_administrators.0.joined_timestamp"),
				),
			},
		},
	})
}

func testAccDataSourceAwsOrganizationsDelegatedAdministratorsConfig(servicePrincipal string) string {
	return testAccAlternateAccountProviderConfig() + fmt.Sprintf(`
data "aws_caller_identity" "delegated" {
  provider = "awsalternate"
}

resource "aws_organizations_delegated_administrator" "test" {
  account_id        = data.aws_caller_identity.delegated.account_id
  service_principal = %[1]q
}

data "aws_organizations_delegated_administrators" "test" {
  service_principal = aws_organizations_delegated_administrator.test.service_principal
}
`, servicePrincipal)
}
