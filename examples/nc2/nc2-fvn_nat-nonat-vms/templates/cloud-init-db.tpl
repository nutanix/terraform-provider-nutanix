#cloud-config
password: "Password1!"
chpasswd: { expire: False }
ssh_pwauth: True
packages:
  - net-tools

users:
  - default
  - name: ubuntu
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    ssh_authorized_keys:
      - ${ssh_public_key}
