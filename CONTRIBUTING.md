## Contributing to the provider

Thank you for your interest in contributing to the Nutanix provider. We welcome your contributions. Here you'll find information to help you get started with provider development.


## Cloning the Project

First, you will want to clone the repository into your working directory:

```shell
git clone git@github.com:nutanix/terraform-provider-nutanix.git
```

## Running the Build

After the clone has been completed, you can enter the provider directory and build the provider.

```shell
cd terraform-provider-nutanix
make build
```

## Developing the Provider

**NOTE:** Before you start work on a feature, please make sure to check the [issue tracker](https://github.com/nutanix/terraform-provider-nutanix/issues) and existing [pull requests](https://github.com/nutanix/terraform-provider-nutanix/pulls) to ensure that work is not being duplicated. For further clarification, you can also ask in a new issue.


If you wish to work on the provider, you'll first need [Go][go-website] installed on your machine.

[go-website]: https://golang.org/
[gopath]: http://golang.org/doc/code.html#GOPATH


## Building Provider

### Installing the Local Plugin

*Note:* manual provider installation is needed only for manual testing of custom
built Nutanix provider plugin.

Manual installation process differs depending on Terraform version.
Run `terraform version` command to determine version of your Terraform installation.

1. Create `/registry.terraform.io/nutanixtemp/nutanix/1.99.99/darwin_amd64/` directories
under:

   * `~/.terraform.d/plugins` (Mac and Linux)

   ```sh
    mkdir -p ~/.terraform.d/plugins/registry.terraform.io/nutanixtemp/nutanix/1.99.99/darwin_amd64/
   ```

2. Build the **binary file**. 
    ```sh 
    go build -o bin/terraform-provider-nutanix_macosx-v1.99.99
    ```

3. Copy Nutanix provider **binary file** there.

   ```sh
    cp bin/terraform-provider-nutanix_macosx-v1.99.99 ~/.terraform.d/plugins/registry.terraform.io/nutanixtemp/nutanix/1.99.99/darwin_amd64/terraform-provider-nutanix_v1.99.99
    cp bin/terraform-provider-nutanix_macosx-v1.99.99 ~/.terraform.d/plugins/terraform-provider-nutanix_v1.99.99
   ```

4. In every Terraform template directory that uses Nutanix provider, ship below
 `terraform.tf` file *(in addition to other Terraform files)*

   ```hcl
   terraform {
     required_providers {
       nutanix = {
         source = "nutanixtemp/nutanix"
         version = "1.99.99"
       }
     }
   }
   ```

5. **Done!**

   Local Nutanix provider plugin will be used after `terraform init`
   command execution in Terraform template directory


### Running tests of provider

For running unit tests:
```sh
make test
```

For running integration tests:

1. Add environment variables for setup related details:
```ssh
export NUTANIX_USERNAME="<username>"
export NUTANIX_PASSWORD="<password>"
export NUTANIX_INSECURE=true
export NUTANIX_PORT=9440
export NUTANIX_ENDPOINT="<pc-ip>"
export NUTANIX_STORAGE_CONTAINER="<storage-container-uuid-for-vm-tests>"
export FOUNDATION_ENDPOINT="<foundation-vm-ip-for-foundation-related-tests>"
export FOUNDATION_PORT=8000
export NOS_IMAGE_TEST_URL="<test-image-url>"
export NDB_ENDPOINT="<ndb-ip>"
export NDB_USERNAME="<username>"
export NDB_PASSWORD="<password>"
```

2. Some tests need setup related constants for resource creation. So add/replace details in test_config.json (for pc tests) and test_foundation_config.json (for foundation and foundation central tests and NDB)

3. To run all tests:
```ssh
make testacc
```

4. To run specific tests:
```ssh 
export TESTARGS='-run=TestAccNutanixPbr_WithSourceExternalDestinationNetwork'
make testacc
```

5. To run collection of tests:
``` ssh
export TESTARGS='-run=TestAccNutanixPbr*'
make testacc
```

### Common Issues using the development binary.

Terraform download the released binary instead developent one.

Just follow this steps to get the development binary:

1. Copy the development terraform binary in the root folder of the project (i.e. where your main.tf is), this should be named `terraform-provider-nutanix`
2. Remove the entire “.terraform” directory.
    ```sh
    rm -rf .terraform/
    ```

3. Run the following command in the same folder where you have copied the development terraform binary.
    ```sh
    terraform init -upgrade
    terraform providers -version
    ```

4. You should see version as “nutanix (unversioned)”
5. Then run your main.tf


## Step to raise a Pull Request

1. Create a github issue with following details.
  * **Title** should contain one of the following
    - [Feat] Develop terraform resource for \<api_name>
    - [Imprv] Modify terraform resource to support \<new_functionality>
    - [Bug] Fix \<summary of issue> bug in \<terraform_resource_name>

  * **Template** for raising issue is already defined. Refer below for creating any issue. 
    - [Issue:Bug Report](https://github.com/nutanix/terraform-provider-nutanix/issues/new?assignees=&labels=&projects=&template=bug_report.md&title=)
    - [Issue:Feature Request](https://github.com/nutanix/terraform-provider-nutanix/issues/new?assignees=&labels=&projects=&template=feature_request.md&title=)

2. Create one of the following git branch from `master` branch. Use `issue#<id>`.
 * `feat/<resource_name>_issue#<id>`
 * `imprv/issue#<id>`
 * `bug/issue#<id>`

3. Please use code comments on the Pull Request. 

4. If a reviewer commented on your PR or asked for changes, please mark the discussion as resolved after you make suggested changes in it. PRs with unresolved issues should not be merged. 
4. Tests are mandatory for each PRs except documentation changes. 
5. Ensure 85% code coverage on the pull request. Pull request with less than 85% coverage will be rejected.
6. Link the pull requested with the associated issue. 
6. Once PR is merged, close the issue. 

## Additional Resources

We've got a handful of resources outside of this repository that will help users understand the interactions between terraform and Nutanix

* YouTube
  _ Overview Video: [](https://www.youtube.com/watch?v=V8_Lu1mxV6g)
  _ Working with images: [](https://www.youtube.com/watch?v=IW0eQevZ73I)
* Nutanix GitHub
  _ [](https://github.com/nutanix/terraform-provider-nutanix)
  _ Private repo until code goes upstream
* Jon’s GitHub
  _ [](https://github.com/JonKohler/ThisOldCloud/tree/master/Terraform-Nutanix)
  _ Contains sample TF’s and PDFs from the youtube videos
* Slack channel \* User community slack channel is available on nutanix.slack.com. Email terraform@nutanix.com to gain entry.

# Nutanix Contributor License Agreement

By submitting a pull request or otherwise contributing to the project, you agree to the following terms and conditions.  You reserve all right and title in your contributions.  

## Grant of License 
You hereby grant Nutanix and to recipients of software distributed by Nutanix, a license to your contributions under the same license as the project.

## Representations
You represent that your contributions are your original creation, and that you are legally entitled to grant the above license.  If your contributions include other third party code, you will include complete details on any third party licenses or restrictions associated with your contributions.

## Notifications
You will notify Nutanix if you become aware that the above representations are inaccurate.  
