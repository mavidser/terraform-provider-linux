# linux_folder

Manages folders and their attributes.

-> Make sure that the user has permissions to the folders being created.

-> If using the provider with a non-sudoer user, allow NOPASSWD sudo access to these commands - `chown` and `chmod`.

## Example Usage

```hcl
resource "linux_user" "testuser" {
  name = "testuser"
  uid = 1024
}

resource "linux_folder" "testfolder" {
  path = "/etc/testfolder"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 777
}
```

## Argument Reference

The following arguments are supported:

- `path` - (Required, string) Absolute path of the folder.
- `owner` - (Optional, string) Owners of the folder, in `user:group` format.
- `permissions` - (Optional, int) Octal permissions of the folder.
