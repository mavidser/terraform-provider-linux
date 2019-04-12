# Terraform linux provider

Just a basic provider to manage linux users/groups on a linux system.
Doesn't support advanced user configuration options (only custom UIDs and GIDs for now).

TODO: Cleanup code, add tests, add more user configuration options

### Sample config

```
provider "linux" {
  host = "192.168.1.2"
  user = "user"
}

resource "linux_group" "testgroup" {
  name = "testgroup"
  system = false
}

resource "linux_user" "testuser1" {
  name = "testuser1"
  gid = "${linux_group.testgroup.gid}"
  system = false
}

resource "linux_user" "testuser2" {
  name = "testuser2"
  gid = "${linux_group.testgroup.gid}"
  system = false
}
```

### Configuration options

#### Provider
- **user**
- **host**
- **port**
- **password**
- **private_key**

#### Group resource
- **name** - Required - int
- **gid** - Optional - int
- **system** - Optional - int

#### User resource
- **name** - Required - int
- **uid** - Optional - int
- **gid** - Optional - int
- **system** - Optional - int
