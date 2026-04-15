#!/bin/bash

# Prism Central credentials
export TF_VAR_NUTANIX_USERNAME="<username>"
export TF_VAR_NUTANIX_PASSWORD="<password>"
export TF_VAR_NUTANIX_ENDPOINT="<your-prism-central-ip>"
# If you are using a self-signed certificate, set this to true
export TF_VAR_NUTANIX_INSECURE="true"
export TF_VAR_NUTANIX_PORT="9440"

# SSH public key
export TF_VAR_SSH_PUBLIC_KEY="ssh-ed25519 AAAAC3NzaC1....."
