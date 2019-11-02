# Linux Provider

This providor is used to manage parts of a typical linux system. Eg - users, groups, files, etc.

## Example Usage

```hcl
provider "linux" {
  host = "192.168.1.128"
  user = "root"
}

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

- `host` - (Required) The host to ssh into.
- `port` - (Optional) The ssh port. Defaults to "22".
- `user` - (Required) The username to ssh with.
- `private_key` - (Optional) The location of the private key, if used for authentication. Defaults to `$HOME/.ssh/id_rsa`.
- `password` - (Optional) The password, if used for authentication.
- `use_sudo` - (Optional) Do certain commands need to be prefixed with sudo? Defaults to true if user is "root", else false.

-> For encrypted private keys, use `ssh-agent` to allow connection.