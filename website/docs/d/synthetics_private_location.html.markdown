---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_private_location"
sidebar_current: "docs-newrelic-datasource-synthetics-private-location"
description: |-
  Grabs a Synthetics monitor location by name.
---

# Data Source: newrelic\_synthetics\_private\_location

Use this data source to get information about a specific Synthetics monitor private location in New Relic that already exists.

## Example Usage

```hcl
data "newrelic_synthetics_private_location" "example" {
  name = "My private location"
}

resource "newrelic_synthetics_monitor" "foo" {
  // Reference the private location data source in the monitor resource
  location_private = [data.newrelic_synthetics_monitor_location.example.id]
}
```

-> This data source doesn't work for `scripted_api`, `scripted_browser` and `step` monitors which works with the latest script runtime.

```hcl
data "newrelic_synthetics_private_location" "example" {
  name = "My private location"
}
resource "newrelic_synthetics_monitor" "foo" {
  // Reference the private location data source in the monitor resource
  location_private { guid = data.newrelic_synthetics_private_location.example.id }
}
```

## Argument Reference

The following arguments are supported:


* `name` - (Required) The name of the Synthetics monitor private location.
