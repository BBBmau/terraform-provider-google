// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package compute_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccComputeImage_imageBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeImageDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeImage_imageBasicExample(context),
			},
			{
				ResourceName:            "google_compute_image.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_encryption_key.0.raw_key", "image_encryption_key.0.rsa_encrypted_key", "labels", "raw_disk", "source_disk", "source_disk_encryption_key", "source_image", "source_image_encryption_key", "source_snapshot", "source_snapshot_encryption_key", "terraform_labels"},
			},
		},
	})
}

func testAccComputeImage_imageBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_image" "debian" {
  family  = "debian-12"
  project = "debian-cloud"
}

resource "google_compute_disk" "persistent" {
  name  = "tf-test-example-disk%{random_suffix}"
  image = data.google_compute_image.debian.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_image" "example" {
  name = "tf-test-example-image%{random_suffix}"

  source_disk = google_compute_disk.persistent.id
}
`, context)
}

func TestAccComputeImage_imageGuestOsExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeImageDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeImage_imageGuestOsExample(context),
			},
			{
				ResourceName:            "google_compute_image.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_encryption_key.0.raw_key", "image_encryption_key.0.rsa_encrypted_key", "labels", "raw_disk", "source_disk", "source_disk_encryption_key", "source_image", "source_image_encryption_key", "source_snapshot", "source_snapshot_encryption_key", "terraform_labels"},
			},
		},
	})
}

func testAccComputeImage_imageGuestOsExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_image" "debian" {
  family  = "debian-12"
  project = "debian-cloud"
}

resource "google_compute_disk" "persistent" {
  name  = "tf-test-example-disk%{random_suffix}"
  image = data.google_compute_image.debian.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_image" "example" {
  name = "tf-test-example-image%{random_suffix}"

  source_disk = google_compute_disk.persistent.id

  guest_os_features {
    type = "UEFI_COMPATIBLE"
  }

  guest_os_features {
    type = "VIRTIO_SCSI_MULTIQUEUE"
  }

  guest_os_features {
    type = "GVNIC"
  }

  guest_os_features {
    type = "SEV_CAPABLE"
  }

  guest_os_features {
    type = "SEV_LIVE_MIGRATABLE_V2"
  }
}
`, context)
}

func TestAccComputeImage_imageBasicStorageLocationExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeImageDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeImage_imageBasicStorageLocationExample(context),
			},
			{
				ResourceName:            "google_compute_image.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_encryption_key.0.raw_key", "image_encryption_key.0.rsa_encrypted_key", "labels", "raw_disk", "source_disk", "source_disk_encryption_key", "source_image", "source_image_encryption_key", "source_snapshot", "source_snapshot_encryption_key", "terraform_labels"},
			},
		},
	})
}

func testAccComputeImage_imageBasicStorageLocationExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_compute_image" "debian" {
  family  = "debian-12"
  project = "debian-cloud"
}

resource "google_compute_disk" "persistent" {
  name  = "tf-test-example-disk%{random_suffix}"
  image = data.google_compute_image.debian.self_link
  size  = 10
  type  = "pd-ssd"
  zone  = "us-central1-a"
}

resource "google_compute_image" "example" {
  name = "tf-test-example-sl-image%{random_suffix}"

  source_disk = google_compute_disk.persistent.id
  storage_locations = ["us-central1"]
}
`, context)
}

func testAccCheckComputeImageDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_image" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/global/images/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ComputeImage still exists at %s", url)
			}
		}

		return nil
	}
}
