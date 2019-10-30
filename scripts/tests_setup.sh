#!/bin/bash
set -e

docker run -d -p 5001:22 --name tf_linux_test_sshd rastasheep/ubuntu-sshd:18.04
export TF_LINUX_SSH_USER=root
export TF_LINUX_SSH_HOST=127.0.0.1
export TF_LINUX_SSH_PORT=5001
export TF_LINUX_SSH_PASSWORD=root