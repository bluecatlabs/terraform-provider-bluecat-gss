# Introduction
---
This is a Terraform provider for the BlueCat Global Server Selector Adaptive Application.


# Requirements
---
- Go version: 1.14
- Terraform version: 0.13.1

# Compile the source code as a Terraform provider
---
- Go to inside of the project directory `cd terraform`
- Compile the project: `go build`

# Running
---
## 1. Preparing the configuration:
---
### Application configuration:
This provider checks for and loads a file "app.yml" from the current working directory, if it exists. The app.yml contains logging configuration for the provider.
Sample "app.yml" configuration file, showing the default values used if no configuration file is present:
```
logging:
  level: warn
  file: provider_bluecat.log
```

### Provider configuration:
Create a file `main.tf` with the following
``` 
provider "bluecatgss" {
  server = "192.0.2.1"
  api_version = "1"
  transport = "https"
  port = "443"
  username = "api_user"
  password = "plain_text_password"
} 
```
server: IP address of the BlueCat Gateway instance running BlueCat Global Server Selector (GSS)
api_version: API version of the GSS application, always 1
transport: http or https
port: Port used to connect to BlueCat Gateway
username: Address Manager API user, used to authenticate to BlueCat Gateway
password: Password used to authenticate to BlueCat Gateway

## 2. Example:
---
```
terraform {
  required_providers {
    bluecatgss = {
      version = ">= 1.0.0"
      source = "terraform-repo.example.com/bluecatlabs/bluecat-gss"
    }
  }
}

provider "bluecatgss" {
  server = "192.0.2.10"
  api_version = "1"
  transport = "http"
  port = "8082"
  username = "gateway"
  password = "123456"
} 

# search order
resource "bluecatgss_search_order" "my_search_order" {
  name = "my-search-order"
  nodes = [	
      {
        name = "Tokyo_Office"  # isolated node
      }
    ]
  links = [
    {
      source = "IO_Office"    # Define 1 direction
      target = "NA_Office"
      cost = 2
      enable_link = true
    },
    {	
      source = "US_Office"   # Define 2 directions
      target = "NA_Office"
      cost = 2
      enable_link = true
    },
    {
      source = "NA_Office"
      target = "US_Office"  
      cost = 2
      enable_link = true
    }
  ]
}

# global application
resource "bluecatgss_application" "my_global_app" {
  configuration = "GSS"
  view = "GSS"
  zone = "example.com"
  absolute_name = "app1-terra"

  health_check {
    type = "TCP"
    check_every = 30
    port = "80"
  }

  search_order = [bluecatgss_search_order.my_search_order.name] # Optional

  fallback = ["10.10.10.11", "10.10.10.12"] # May define CNAME
}

# answers
resource "bluecatgss_answer" "na-server" {
  application_id = bluecatgss_application.my_global_app.id
  addresses = [ "192.168.41.198"]
  region = "NA_Office"
  type = "ip_address"
  name = "na-01"
}

resource "bluecatgss_answer" "us-server" {
  application_id = bluecatgss_application.my_global_app.id
  addresses = [ "www.google.com"]
  region = "US_Office"
  type = "fqdn"
  name = "us-01"
}
```

## 3. Preparing the resources:
---

### Resource GSS Search Order:
```
resource "bluecatgss_search_order" "search_order" {
  name = "my_search_order"
  nodes = [
    {
      name = "US"
    }
  ]
  links = [
    {
      source = "Tokio"
      target = "Denver"
      cost = 10
      enable_link = false
    }
  ]
}
```

The following arguments are supported:

* ```name``` - (Required) name of search order.
* ```nodes``` - (Required) Defined nodes. List any isolated nodes that have no links
* ```links``` - (Required) Specifies the source, target, cost, enable_link to add a link between nodes

* Notes: enable_link is optional and default value is True. This attribute is required in Attributes as Blocks.
* For a standard bidirectional link between two sites, add two links in opposite directions.

### Resource GSS Application:
```
resource "bluecatgss_application" "global_application" {
  configuration = "terraform_demo"
  view = "Internal"
  zone = "example.com"
  absolute_name = "app"
  ttl = 60
  properties = ""
  fallback = ["10.10.10.20"]
  health_check {}
  search_order = [ bluecatgss_search_order.search_order.name ]
}
```
The following arguments are supported:

* ```configuration``` - (Required) The Configuration in Address Manager.
* ```view``` - (Required) The DNS View in Address Manager.
* ```zone``` - (Optional) The absolute name of the DNS zone in which you want to update a Host record.
* ```absolute_name``` - (Required) The name of the Host record. This must be an FQDN if the zone is not provided.
* ```ttl``` - (Optional) The TTL value. Default is -1, indicating the inherited value is used.
* ```properties``` - (Optional) Additional properties to be passed when adding the Host record.
* ```fallback``` - (Required) Specifies a single FQDN or set of IP addresses that are used for the Fallback answers in GSS, and the initial value of the Host record (or CNAME).
* ```health_check``` - (Optional) Specifies the Health Check Configuration: TCP, HTTP_HEAD, CUSTOMIZE, NO_HEALTH_CHECK.

Examples of Health Check configuration for each type of Health Check:

1. TCP Type :
```
health_check {
  type = "TCP"
  check_every = 30
  port = "22"
}
```

2. HTTP Head Type:
```
health_check {
  type = "HTTP_HEAD"
  check_every = 30
  secure_connection = true
  appended_url_path = "/"
  optional_header = ""
  header_value = ""
}
```

3. Customize Type:
```
health_check {
  type = "CUSTOMIZE"
  check_every = 30
  custom_data = {
    "Field 1" = "Data 1"
    "Field 2" = "Data 2"
}
```

4. No Health Check Type:
```
health_check {
  type = "NO_HEALTH_CHECK"
}
```

If no health_check configuration is provided, then the application is configured for no health check.

* ```search_order``` - (Optional) Defines the Search Order configurations to link to this application. For example:
```
search_order = [bluecatgss_search_order.search_order.name]
```

# Resource GSS Answer:
```
resource "bluecatgss_answer" "answer_region" {
    application_id = bluecatgss_application.global_application.id
    addresses = [ "10.10.10.10", "10.10.10.11"]
    region = "answer_region1"
    name = "demo-region"
    type = "ip_address"
}
```
The following arguments are supported:

* ```application_id``` - (Required) The GSS Application Resource ID. This the object ID of the application Host record or CNAME record in Address Manager.
* ```addresses``` - (Required) Specifies the address(es) of Answer.
* ```region``` - (Required) Specifies the region of Answer.
* ```name``` - (Required) Specifies the name of Answer.
* ```type``` - (Required) Specifies the type of Answer. Value is either *ip_address* or *fqdn*

**Note** If the region of an Answer resource is changed. Any existing answer will be destroyed, and a new answer created.

## 4. Executing the provider:
---
### Initialize the provider:

In case of you're using a local build of the provider, you need to prepare the directory structure to install the provider as described below. Otherwise, just run `terraform init`.
- Create the directory to store the providers: <HOME_DIR>/providers
- Create the provider structure for the provider under the directory at step 1: <HOSTNAME>/<NAMESPACE>/<TYPE>/<VERSION>/<PLATFORM>/<PROVIDER_BINARY>. For example: example.com/bluecatlabs/bluecat-gss/1.0.0/windows_amd64/terraform-provider-bluecat-gss.exe
- Add a provider block to your configuration:
    ```
    terraform {
      required_providers {
        <TYPE> = {
          version = ">= <VERSION>"
          source = "<HOSTNAME>/<NAMESPACE>/<TYPE>"
        }
      }
    }
    ```
  For example:
    ```
    terraform {
      required_providers {
        bluecatgss = {
          version = ">= 1.0.0"
          source = "example.com/bluecatlabs/bluecat-gss"
        }
      }
    }
    ```
- Install your provider: `terraform init -plugin-dir=<HOME_DIR>/providers`

### Plan and review required changes
`terraform plan`

### Add/update resources
`terraform apply`

### Remove managed resources
`terraform destroy`
