# linux_group

Manages groups and their attributes.

-> If using the provider with a non-sudoer user, allow NOPASSWD sudo access to these commands - `groupadd`, `groupmod`, and `groupdel`.

## Example Usage

```hcl
resource "linux_group" "testgroup" {
  name = "testgroup"
  gid = 1048
}

resource "linux_user" "testuser" {
  name = "testuser"
  uid = 1024
  gid = linux_group.testgroup.gid
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required, string) Name of the group.
- `gid` - (Optional, int) gid to set of the group.
- `system` - (Optional, bool) If GID is not supplied, this attribute is factored in while generating the GID. Defaults to false.

## Attribute Reference

The following attributes are exported:

- `gid` - If not supplied, the generated GID.