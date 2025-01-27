echo "Starting build Mac"
export GOOS="linux"
export GOARCH="amd64"
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/terraform-providers/nutanix/1.99.99/linux_amd64/
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/nutanixtemp/nutanix/1.99.99/linux_amd64/
go build -o bin/terraform-provider-nutanix_macosx-v1.99.99
unset GOOS
unset GOARCH
echo "Finished build... Starting copy"
cp bin/terraform-provider-nutanix_macosx-v1.99.99 ~/.terraform.d/plugins/registry.terraform.io/terraform-providers/nutanix/1.99.99/linux_amd64/terraform-provider-nutanix_v1.99.99
cp bin/terraform-provider-nutanix_macosx-v1.99.99 ~/.terraform.d/plugins/registry.terraform.io/nutanixtemp/nutanix/1.99.99/linux_amd64/terraform-provider-nutanix_v1.99.99
cp bin/terraform-provider-nutanix_macosx-v1.99.99 ~/.terraform.d/plugins/terraform-provider-nutanix_v1.99.99
echo "deleting terraform dependency lock file"
rm ./temp/.terraform.lock.hcl
rm -rf ./temp/.terraform
rm ./temp/terraform.tfstate
rm ./temp/terraform.tfstate.backup
echo "Done"
