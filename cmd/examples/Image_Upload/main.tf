provider "nutanix" { 
    username = "admin"
    password = "Nutanix/1234"
    endpoint = "10.5.80.30"
    insecure = true
}

resource "nutanix_image" "centos7-iso" {
    name = "CentOS7-ISO"
    source_uri = "http://10.7.1.7/data1/ISOs/CentOS-7-x86_64-Minimal-1503-01.iso" 
    checksum_algorithm = "SHA_256"
    checksum_value = "a9e4e0018c98520002cd7cf506e980e66e31f7ada70b8fc9caa4f4290b019f4f"
}
resource "nutanix_image" "centos-base-image" {
    name = "Centos7-Base-Image"
    source_uri = "http://10.7.1.7/data1/AHVUVMImages/Centos7-Base.qcow2"
}
/*
resource "nutanix_image" "centos7-iso-File" {
    name = "CentOS7-ISO-File"
    source_uri = "file://CentOS-7-x86_64-Minimal-1503-01.iso" 
    checksum_algorithm = "SHA_256"
    checksum_value = "a9e4e0018c98520002cd7cf506e980e66e31f7ada70b8fc9caa4f4290b019f4f"
}

resource "nutanix_image" "centos-base-image-File" {
    name = "Centos7-Base-Image-File"
    source_uri = "file://Centos7-Base.qcow2"
}
*/
