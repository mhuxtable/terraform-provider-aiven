// Copyright (c) 2017 jelmersnoeck
// Copyright (c) 2018-2021 Aiven, Helsinki, Finland. https://aiven.io/
package aiven

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Kafka service tests
func TestAccAivenService_kafka(t *testing.T) {
	resourceName := "aiven_kafka.bar"
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAivenServiceResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaServiceResource(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAivenServiceCommonAttributes("data.aiven_kafka.service"),
					testAccCheckAivenServiceKafkaAttributes("data.aiven_kafka.service"),
					resource.TestCheckResourceAttr(resourceName, "service_name", fmt.Sprintf("test-acc-sr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "state", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "project", os.Getenv("AIVEN_PROJECT_NAME")),
					resource.TestCheckResourceAttr(resourceName, "service_type", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "cloud_name", "google-europe-west1"),
					resource.TestCheckResourceAttr(resourceName, "maintenance_window_dow", "monday"),
					resource.TestCheckResourceAttr(resourceName, "maintenance_window_time", "10:00:00"),
					resource.TestCheckResourceAttr(resourceName, "state", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "termination_protection", "false"),
				),
			},
		},
	})
}

func testAccKafkaServiceResource(name string) string {
	return fmt.Sprintf(`
		data "aiven_project" "foo" {
			project = "%s"
		}
		
		resource "aiven_kafka" "bar" {
			project = data.aiven_project.foo.project
			cloud_name = "google-europe-west1"
			plan = "business-4"
			service_name = "test-acc-sr-%s"
			maintenance_window_dow = "monday"
			maintenance_window_time = "10:00:00"
			
			kafka_user_config {
				kafka_rest = true
				kafka_connect = true
				schema_registry = true

				kafka {
					group_max_session_timeout_ms = 70000
					log_retention_bytes = 1000000000
				}

				public_access {
					kafka_rest = true
					kafka_connect = true
				}
			}
		}
		
		data "aiven_kafka" "service" {
			service_name = aiven_kafka.bar.service_name
			project = aiven_kafka.bar.project

			depends_on = [aiven_kafka.bar]
		}
		`, os.Getenv("AIVEN_PROJECT_NAME"), name)
}

func testAccCheckAivenServiceKafkaAttributes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["kafka_user_config.0.kafka_connect"] != "true" {
			return fmt.Errorf("expected to get a correct kafka_connect from Aiven")
		}

		if a["kafka_user_config.0.kafka_rest"] != "true" {
			return fmt.Errorf("expected to get a correct kafka_rest from Aiven")
		}

		if a["kafka_user_config.0.public_access.0.kafka_connect"] != "true" {
			return fmt.Errorf("expected to get a correct public_access.kafka_connect from Aiven")
		}

		if a["kafka_user_config.0.public_access.0.kafka_rest"] != "true" {
			return fmt.Errorf("expected to get a correct public_access.kafka_rest from Aiven")
		}

		if a["kafka_user_config.0.public_access.0.kafka"] != "" {
			return fmt.Errorf("expected to get a correct public_access.kafka from Aiven")
		}

		if a["kafka_user_config.0.public_access.0.prometheus"] != "" {
			return fmt.Errorf("expected to get a correct public_access.prometheus from Aiven")
		}

		if a["kafka.0.connect_uri"] == "" {
			return fmt.Errorf("expected to get a connect_uri from Aiven")
		}

		if a["kafka.0.rest_uri"] == "" {
			return fmt.Errorf("expected to get a rest_uri from Aiven")
		}

		if a["kafka.0.schema_registry_uri"] == "" {
			return fmt.Errorf("expected to get a schema_registry_uri from Aiven")
		}

		if a["kafka.0.access_key"] == "" {
			return fmt.Errorf("expected to get an access_key from Aiven")
		}

		if a["kafka.0.access_cert"] == "" {
			return fmt.Errorf("expected to get an access_cert from Aiven")
		}

		return nil
	}
}
