#cloud-config
user: nutanix
password: <password> 
chpasswd: {expire: False}
sudo: ALL=(ALL) NOPASSWD:ALL
ssh_pwauth: True
fqdn: ${hostname}.test.local
hostname: ${hostname}

apt_upgrade: true
packages:
   - python-is-python3
#   - python3-pip


runcmd:
  - [mkdir, /mnt/nutanix]
  - [mount, /dev/sr1, /mnt/nutanix]
  - [python, /mnt/nutanix/installer/linux/install_ngt.py]