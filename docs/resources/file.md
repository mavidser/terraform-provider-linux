# linux_file

Manages files and their attributes.

-> Make sure that the user has permissions to the files being created.

## Example Usage

```hcl
resource "linux_user" "testuser" {
  name = "testuser"
  uid = 1024
}

resource "linux_file" "testfile" {
  path = "/etc/testfile"
  content = "testcontent"
  owner = "${linux_user.testuser.name}:${linux_user.testuser.name}"
  permissions = 777
}
```

## Argument Reference

The following arguments are supported:

- `path` - (Required, string) Absolute path of the file.
- `owner` - (Optional, string) Owners of the file, in `user:group` format.
- `permissions` - (Optional, int) Octal permissions of the file.
- `content` - (Optional, string) Content of the file.
