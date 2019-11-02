# linux_user

Manages users and their attributes.

-> If using the provider with a non-sudoer user, allow NOPASSWD sudo access to these commands - `useradd`, `usermod`, and `userdel`.

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

- `name` - (Required, string) Name of the user.
- `uid` - (Optional, int) UID to set of the user.
- `gid` - (Optional, int) GID to set of the user.
- `system` - (Optional, bool) If UID is not supplied, this attribute is factored in while generating the GID. Defaults to false.

## Attribute Reference

The following attributes are exported:

- `uid` - If not supplied, the generated uid.
- `gid` - If not supplied, the generated uid.