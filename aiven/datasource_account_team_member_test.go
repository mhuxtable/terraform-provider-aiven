// Copyright (c) 2017 jelmersnoeck
// Copyright (c) 2018-2021 Aiven, Helsinki, Finland. https://aiven.io/
package aiven

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAivenAccountTeamMemberDataSource_basic(t *testing.T) {
	datasourceName := "data.aiven_account_team_member.member"
	resourceName := "aiven_account_team_member.foo"
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountTeamMemberResource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "account_id", resourceName, "account_id"),
					resource.TestCheckResourceAttrPair(datasourceName, "user_email", resourceName, "user_email"),
					resource.TestCheckResourceAttrPair(datasourceName, "team_id", resourceName, "team_id"),
					resource.TestCheckResourceAttrPair(datasourceName, "create_time", resourceName, "create_time"),
					resource.TestCheckResourceAttrPair(datasourceName, "accepted", resourceName, "accepted"),
				),
			},
		},
	})
}
