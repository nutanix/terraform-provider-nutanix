## 2.1.1 (April 08, 2025)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v2.1.0...v2.1.1)

**Fixed bugs:**
- Error Creating NoNat Overlay External Network [\#820](https://github.com/nutanix/terraform-provider-nutanix/issues/820)

## 2.1.0 (April 07, 2025)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v2.0.0...v2.1.0)

All new features are v4 SDKs based. 

**New Feature:**
- Prism [\#815](https://github.com/nutanix/terraform-provider-nutanix/issues/815)
    - CRUD for Backup Target
    - CRD for Restore Source
    - Resource to deploy pc
    - Resource to restore pc
    - Restore to unregister pc
    - Datasource to list pcs
    - Datasource to fetch pc details
    - Datasource to list restorable pcs
    - Datasource to list pc restore points
    - Datasource to list pc restore points details
    - Datasource to list backup targets
    - Datasource to fetch backup target
    - Datasource to fetch restore source

- Data Protection [\#816](https://github.com/nutanix/terraform-provider-nutanix/issues/816)
    - CRUD for protection policies
    - Resource to promote protected resource
    - Resource to restore protected resource
    - Datasource to fetch protected resource
    - Datasource to list protection policies
    - Datasource to fetch protection policy

- LCM [\#814](https://github.com/nutanix/terraform-provider-nutanix/issues/814)
   - Resource for Perform Inventory
   - Resource for Upgrade PreChecks
   - Resource for Upgrade of entities
   - Resource for LCM Config
   - Datasource to fetch LCM status
   - Datasource to fetch entities
   - Datasource to fetch specific entity
   - Datasource to fetch LCM Config

**Fixed bugs:**
- Fix to add additional disk for virtual machine during update [\#750](https://github.com/nutanix/terraform-provider-nutanix/issues/750)
- Fix to create a virtual machine with UEFI [\#725](https://github.com/nutanix/terraform-provider-nutanix/issues/725)
- Fix Nutanix image v2 example in documentation [\#738](https://github.com/nutanix/terraform-provider-nutanix/issues/738)
- Fix datasource to expand cluster v2 [\#722](https://github.com/nutanix/terraform-provider-nutanix/issues/722)
- Support import in nutanix_virtual_machine_v2 Resource [\#745](https://github.com/nutanix/terraform-provider-nutanix/issues/745)
- Fix documentation and examples for existing features [\#755](https://github.com/nutanix/terraform-provider-nutanix/pull/755)


## 2.0.0 (January 07, 2025)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.9.5...v2.0.0)

All new features are v4 SDKs based. 

**New Feature:**
- Cluster Management [\#704](https://github.com/nutanix/terraform-provider-nutanix/issues/704)
    - CRUD for clusters
    - Resource to discover unconfigured nodes
    - Resource to fetch network configuration of unconfigured nodes
    - Resource to add/remove node from cluster
    - Resource for PC registration
    - Datasource for Hosts Info

- Storage Containers [\#705](https://github.com/nutanix/terraform-provider-nutanix/issues/705)
    - CRUD for storage containers
    - datasource for storage stats info

- Networks [\#706](https://github.com/nutanix/terraform-provider-nutanix/issues/706)
    - CRUD for Subnets
    - CRUD for VPCs
    - CRUD for PBRs
    - CRUD for Floating IPs
    - CRUD for Static Routes
    - CRUD for Services Groups
    - CRUD for Address Groups
    - CRUD for Network Security Rules

- IAM [\#707](https://github.com/nutanix/terraform-provider-nutanix/issues/707)
    - CRUD for Roles
    - Datasource for Operations( permissions )
    - CRUD for Users
    - CRUD for User Groups
    - CRUD for Authorization Policy ( ACPs )
    - CRUD for SALM Identity Providers
    - CRUD for Directory Services

- Prism [\#711](https://github.com/nutanix/terraform-provider-nutanix/issues/711)
    - CRUD for categories

- VMM [\#708](https://github.com/nutanix/terraform-provider-nutanix/issues/708)
    - CRUD for Virtual Machine
    - CRUD for Images
    - CRUD for Image Placement Policy
    - Resource to Clone a VM
    - Resource for NGT Installation / Uninstallation
    - Resource for NGT upgrade
    - Resource for Template deploy
    - Resource for Template Guest OS Actions
    - CRUD for Template
    - Resource for CDROM inject/eject
    - Resource for VM revert from recovery point
    - Resource for VM guest customisation update
    - Resource for VM Network Device Assign/Remove IP
    - Resource for VM Network Device Migrate
    - Resource for VM Shutdown Actions

- Volumes [\#709](https://github.com/nutanix/terraform-provider-nutanix/issues/709)
    - CRUD for Volume Groups
    - Resource to Attach/Deattach VG to VM
    - Resource to Attach/Deattach VG to ISCSI Client
    - CRUD for Volume Group vDisks
    - Resource to attach/deattach categories from VG

- Data Protection [\#710](https://github.com/nutanix/terraform-provider-nutanix/issues/710)
    - CRUD for Recovery Point
    - CRUD for Replicate Recovery Point
    - Resource to Restore VM from Recovery Point
    - Datasource for VM Recovery Point Info



## 1.9.5 (January 16, 2024)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/feat/1.9.4...feat/1.9.5)

**Merged pull request:**
- Fixing the error for NDB datasource for databases if creds are invalid. [\#668](https://github.com/nutanix/terraform-provider-nutanix/pull/668)

**Closed issues:**
- Memory MiB quantity value shows discrepancy in terraform plugin and Prism.. [\#654](https://github.com/nutanix/terraform-provider-nutanix/issues/654)

**Fixed bugs:**
- Datasource nutanix_ndb_databases doesn't return error if creds are invalid. [\#667](https://github.com/nutanix/terraform-provider-nutanix/issues/667)


## 1.9.4 (October 27, 2023)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/feat/1.9.3...feat/1.9.4)

- Feat branch. [\#645](https://github.com/nutanix/terraform-provider-nutanix/pull/645)

**Merged pull request:**
- Change VM name should not require VM PowerOFF. [\#626](https://github.com/nutanix/terraform-provider-nutanix/pull/626)
- Fix: compare bootconfig against previous value. [\#641](https://github.com/nutanix/terraform-provider-nutanix/pull/641)

**Implemented enhancements:**
- Added example to use metadata in nutanix subnets. [\#643](https://github.com/nutanix/terraform-provider-nutanix/pull/643)
- External subnet name/uuid are Optional args not Required. [\#644](https://github.com/nutanix/terraform-provider-nutanix/pull/644)

**Fixed bugs:**
- VM rebooted at every change because of hotPlugChange set to false. [\#640](https://github.com/nutanix/terraform-provider-nutanix/issues/640)
- Changing the VM name forces a reboot. [\#625](https://github.com/nutanix/terraform-provider-nutanix/issues/625)

**Closed issues:**
- Modify Terraform documentation for nutanix_vpc resource. [\#636](https://github.com/nutanix/terraform-provider-nutanix/issues/636)
- Include metadata example for data.nutanix_subnets. [\#590](https://github.com/nutanix/terraform-provider-nutanix/issues/590)


## 1.9.3 (September 7, 2023)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/feat/1.9.2...feat/1.9.3)

**Merged pull request:**
- Setting machine type in updating virtual machines. [\#630](https://github.com/nutanix/terraform-provider-nutanix/pull/630)
- Added examples of role creation using nutanix terraform provider. [\#632](https://github.com/nutanix/terraform-provider-nutanix/pull/632)

**Fixed bugs:**
- Updating gives error: Machine type must be set to Q35 for secure boot. [\#622](https://github.com/nutanix/terraform-provider-nutanix/issues/622)
- Machine type must be set to Q35 for secure boot. [\#494](https://github.com/nutanix/terraform-provider-nutanix/issues/494)

**Closed issues:**
- Add support documentation in terraform. [\#611](https://github.com/nutanix/terraform-provider-nutanix/issues/611)

**Closed pull request:**
-  Fix Secure boot VMs when doing updates. [\#496](https://github.com/nutanix/terraform-provider-nutanix/pull/496)


## 1.9.2 (July 21, 2023)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/feat/1.9.1...feat/1.9.2)

**Merged pull request:**
- Fixing terraform doc link for detailed references. [\#618](https://github.com/nutanix/terraform-provider-nutanix/pull/618)

**Closed issues:**
- All "see detailed information" at the bottom of the doc pages are broken. [\#617](https://github.com/nutanix/terraform-provider-nutanix/issues/617)
- Subnet documentation missing attributes. [\#616](https://github.com/nutanix/terraform-provider-nutanix/issues/616)
- Unsupported attribute for Volume Groups should be removed from VM resource documentation. [\#613](https://github.com/nutanix/terraform-provider-nutanix/issues/613)
- Typo in documentation guest_customization_sysrep_custom_key_values. [\#495](https://github.com/nutanix/terraform-provider-nutanix/issues/495)

## 1.9.1 (June 14, 2023)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/feat/1.9.0...feat/1.9.1)

**Implemented enhancements:**
- Handling version in ReadContext for "nutanix_karbon_cluster" resource. [\#606](https://github.com/nutanix/terraform-provider-nutanix/pull/606)
- Changed Port field in "resource_karbon_private_registry" to Optional. [\#607](https://github.com/nutanix/terraform-provider-nutanix/pull/607)
- NDB clusters endpoint to have support for Active Directory. [\#608](https://github.com/nutanix/terraform-provider-nutanix/pull/608)

**Closed issues:**
- Add support documentation in terraform. [\#611](https://github.com/nutanix/terraform-provider-nutanix/issues/611)


## 1.9.0 (May 26, 2023)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/feat/1.8.1...feat/1.9.0)

**New Feature:**
- Add support for new Karbon features. [\#290](https://github.com/nutanix/terraform-provider-nutanix/issues/290)
    
    New Resource :
    - nutanix_karbon_worker_nodepool

**Implemented enhancements:**
- Adding timeouts in "nutanix_karbon_cluster" resource. [\#563](https://github.com/nutanix/terraform-provider-nutanix/pull/563)
- Vlan with 0 vlan_id should be supported in subnet resource. [\#562](https://github.com/nutanix/terraform-provider-nutanix/pull/562)
- Adding contributing doc and code of conduct in provider. [\#603](https://github.com/nutanix/terraform-provider-nutanix/pull/603)
- Schema Validation for NDB database provision when profiles are required or optional. [\#591](https://github.com/nutanix/terraform-provider-nutanix/issues/591)

**Fixed bugs:**
- Intermittent "context deadline exceeded" errors on "nutanix_karbon_cluster" resource. [\#544](https://github.com/nutanix/terraform-provider-nutanix/issues/544)
- Resource "nutanix_subnet" fails when creating a managed IPAM network using a VLAN that overlaps with existing network. [\#543](https://github.com/nutanix/terraform-provider-nutanix/issues/543)
- In NDB database resource, Required profile values for provisioning a database with registered dbserver or with new dbserver should be properly listed. [#\604](https://github.com/nutanix/terraform-provider-nutanix/issues/604)

**Closed issues:**
- Typo in documentation guest_customization_sysrep_custom_key_values. [\#495](https://github.com/nutanix/terraform-provider-nutanix/issues/495)
- Documentation about subnet_type for nutanix_subnet is missing. [\#506](https://github.com/nutanix/terraform-provider-nutanix/issues/506)
- parent_reference misspelled. [\#507](https://github.com/nutanix/terraform-provider-nutanix/issues/507)
- availability_zone_reference not returning in nutanix_clusters. [\#573](https://github.com/nutanix/terraform-provider-nutanix/issues/573)

**Merged pull requests:**
- Add information about how to import virtual machine to state. [\#500](https://github.com/nutanix/terraform-provider-nutanix/pull/500)
- Removed UUID field from documentation of nutanix address group. [\#462](https://github.com/nutanix/terraform-provider-nutanix/pull/462)


## 1.8.1 (April 18, 2023)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/feat/1.8.0-ga...feat/1.8.1)

**New Feature:**
- Added a datasource for NDB network available ips. [\#570](https://github.com/nutanix/terraform-provider-nutanix/pull/570)

    New Data Source :
    -  nutanix_ndb_network_available_ips

**Implemented enhancements:**
- Handling the case if ndb cluster id is given wrong. [\#571](https://github.com/nutanix/terraform-provider-nutanix/pull/571)
- Support for imports on ndb day-2 actions. [\#561](https://github.com/nutanix/terraform-provider-nutanix/pull/561)
- Changing the provider version in ndb examples from beta to ga. [\#580](https://github.com/nutanix/terraform-provider-nutanix/pull/580)

**Fixed bugs:**
- Tf ntnx provider crash for wrong ndb cluster id. [\#567](https://github.com/nutanix/terraform-provider-nutanix/issues/567)
- lcm_config should be Set without plugin crash in clone resource. [\#583](https://github.com/nutanix/terraform-provider-nutanix/issues/583)

**Closed issues:**
- NDB datasource for network available ips. [\#569](https://github.com/nutanix/terraform-provider-nutanix/issues/569) 
- Support for adding import on ndb day2 actions. [\#582](https://github.com/nutanix/terraform-provider-nutanix/issues/582)
- How to configure IP address for Nutanix Packer images. [\#576](https://github.com/nutanix/terraform-provider-nutanix/issues/576)


## 1.8.0 (Feb 23, 2023)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/feat/v1.8.0-beta-2...feat/1.8.0-ga)

**New Features:**
- Feat/1.8.0-ga Release with Nutanix Database Service based resource and datasources [\#553](https://github.com/nutanix/terraform-provider-nutanix/pull/553)

    New Resources:
    - nutanix_ndb_maintenance_window
    - nutanix_ndb_maintenance_task
    - nutanix_ndb_tms_cluster
    - nutanix_ndb_tag
    - nutanix_ndb_network
    - nutanix_ndb_dbserver_vm
    - nutanix_ndb_register_dbserver
    - nutanix_ndb_stretched_vlan
    - nutanix_ndb_clone_refresh
    - nutanix_ndb_cluster

    New Data Sources:
    - nutanix_ndb_maintenance_window
    - nutanix_ndb_maintenance_windows
    - nutanix_ndb_tag
    - nutanix_ndb_tags
    - nutanix_ndb_network
    - nutanix_ndb_networks
    - nutanix_ndb_dbserver
    - nutanix_ndb_dbservers
    

## 1.8.0-beta-2 (Jan 20, 2023)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.8.0-beta.1...v1.8.0-beta.2)

**New Feature:**

- Feat/1.8.0-beta.2 Release with Nutanix Database Service based resource and datasources [\#533] (https://github.com/nutanix/terraform-provider-nutanix/pull/533)

    New Resources:
    - nutanix_ndb_profile
    - nutanix_ndb_sla
    - nutanix_ndb_database_scale
    - nutanix_ndb_database_restore
    - nutanix_ndb_database_snapshot
    - nutanix_ndb_register_database
    - nutanix_ndb_clone
    - nutanix_ndb_log_catchups
    - nutanix_ndb_authorize_dbservers
    - nutanix_ndb_software_version_profile
    - nutanix_ndb_linked_databases

    New Data Sources:
    - nutanix_ndb_snapshot
    - nutanix_ndb_snapshots
    - nutanix_ndb_time_machine
    - nutanix_ndb_time_machines
    - nutanix_ndb_tms_capability
    - nutanix_ndb_clone
    - nutanix_ndb_clones


**Implemented enhancements:**
 - Support for HA instance in nutanix_ndb_database resource. [\#518](https://github.com/nutanix/terraform-provider-nutanix/pull/518)
 - Improving the error when server is unreachable. [\#530](https://github.com/nutanix/terraform-provider-nutanix/pull/530)
 - Fetching of database based on database_type filter [\#513](https://github.com/nutanix/terraform-provider-nutanix/pull/513)
 - Support of Tags and Maintainence Window in provisioning [\#528](https://github.com/nutanix/terraform-provider-nutanix/pull/528) 


## 1.8.0-beta.1 (Oct 12, 2022)

[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.7.1...v1.8.0-beta.1)

**New Feature:**

- Feat/1.8.0-beta.1 Release with Nutanix Database Service based resources and datasources  [\#501] (https://github.com/nutanix/terraform-provider-nutanix/pull/501)

    New Data Sources :
    - nutanix_ndb_cluster
    - nutanix_ndb_clusters
    - nutanix_ndb_database
    - nutanix_ndb_databases
    - nutanix_ndb_profile
    - nutanix_ndb_profiles
    - nutanix_ndb_sla
    - nutanix_ndb_slas

    New Resources :
    - nutanix_ndb_database

## 1.7.1 (August 31, 2022)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.7.0...v1.7.1)

**Implemented enhancements:**

- Support for /projects_internal API in nutanix projects [\#487](https://github.com/nutanix/terraform-provider-nutanix/pull/487)

**Closed Issues:**
- project internal changes [\#488] (https://github.com/nutanix/terraform-provider-nutanix/pull/488)


## 1.7.0 (August 12, 2022)

[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.6.1...v1.7.0)

**New Feature:**
- Added a resource for user groups. [\#477] (https://github.com/nutanix/terraform-provider-nutanix/pull/477)

    New Resource :
    -  nutanix_user_groups

**Merged pull request:**

- Feature user groups [\#477] (https://github.com/nutanix/terraform-provider-nutanix/pull/477) ([abhimutant](https://github.com/abhimutant))
- Fix checksum issue in image resource [\#481] (https://github.com/nutanix/terraform-provider-nutanix/pull/481) ([bhatipradeep](https://github.com/bhatipradeep))
- Fix image resource issues [\#480] (https://github.com/nutanix/terraform-provider-nutanix/pull/480) ([bhatipradeep](https://github.com/bhatipradeep))
- Add field to set time zone for cluster created using foundation [\#484] (https://github.com/nutanix/terraform-provider-nutanix/pull/484) ([bhatipradeep](https://github.com/bhatipradeep))
- fixed the issue and added testcase for sysprep custom key values guest customisation [\#478] (https://github.com/nutanix/terraform-provider-nutanix/pull/478) ([abhimutant](https://github.com/abhimutant))
- Changes in karbon cluster documentation for storage config [\#479] (https://github.com/nutanix/terraform-provider-nutanix/pull/479) ([abhimutant](https://github.com/abhimutant))
- Enable tests for protection rules based datasources and resources [\#366] (https://github.com/nutanix/terraform-provider-nutanix/pull/366) ([bhatipradeep](https://github.com/bhatipradeep))
- Add uuid to address_groups data source [\#472] (https://github.com/nutanix/terraform-provider-nutanix/pull/472) ([mcguppy](https://github.com/mcguppy))


**Fixed bugs:**

- Terraform provider crashes when using guest_customization_sysprep_custom_key_values [\#441] (https://github.com/nutanix/terraform-provider-nutanix/issues/441) 
- Nutanix terraform- Karbon clusters, the storage_class_config is not been displayed. [\#417] (https://github.com/nutanix/terraform-provider-nutanix/issues/417)
- Checksum is not considered while uploading image from local using nutanix_image resource. [\#469] (https://github.com/nutanix/terraform-provider-nutanix/issues/469)
- Not able to update image_type of PC image [\#454] (https://github.com/nutanix/terraform-provider-nutanix/issues/454)
- Image gets deleted from PC if nutanix_image update errors out [\#453] (https://github.com/nutanix/terraform-provider-nutanix/issues/453)

**Closed issues:**

- Support for User Groups [\#475] (https://github.com/nutanix/terraform-provider-nutanix/issues/475)
- Resizing disk identified in plan, but not actually done during apply [\#463] (https://github.com/nutanix/terraform-provider-nutanix/issues/463) 
- uuid of address_groups are not available [\#461] (https://github.com/nutanix/terraform-provider-nutanix/issues/461)
- ntx provider have ENTITY_READ_ERROR when try to recreate a VM deleted manually [\#451] (https://github.com/nutanix/terraform-provider-nutanix/issues/451)
- Allow project definiation by name and not just ID [\#406] (https://github.com/nutanix/terraform-provider-nutanix/issues/406)
- Upload images to Foundation VM from a remote server using nutanix_foundation_image resource [\#455] (https://github.com/nutanix/terraform-provider-nutanix/issues/455)
- Change catagories from list to dictionary for more straightforward use [\#407] (https://github.com/nutanix/terraform-provider-nutanix/issues/407)
- VM creation fails when boot_type = LEGACY [\#304] (https://github.com/nutanix/terraform-provider-nutanix/issues/304)
- Add attributes to set timezone during cluster creation using image_nodes [\#449] (https://github.com/nutanix/terraform-provider-nutanix/issues/449)
- Enable multi-pc tests [\#448] (https://github.com/nutanix/terraform-provider-nutanix/issues/448)


## 1.6.1 (July 7, 2022)

[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.5.0...v1.6.1)

**New Feature:**

- Add Datasource and Resource for Flow Networking [\#473] (https://github.com/nutanix/terraform-provider-nutanix/pull/473)

    New Data Sources :
    - nutanix_floating_ip
    - nutanix_floating_ips
    - nutanix_pbr
    - nutanix_pbrs
    - nutanix_static_routes
    - nutanix_vpc
    - nutanix_vpcs

    New Resources :
    - nutanix_floating_ip
    - nutanix_pbr
    - nutanix_static_routes
    - nutanix_vpc

**Merged pull requests:**

- VPC Feature [\#457](https://github.com/nutanix/terraform-provider-nutanix/pull/457) ([abhimutant](https://github.com/abhimutant))
- Policy Based Routing Feature [\#464](https://github.com/nutanix/terraform-provider-nutanix/pull/464) ([abhimutant](https://github.com/abhimutant))
- Floating IP Feature [\#466](https://github.com/nutanix/terraform-provider-nutanix/pull/466) ([abhimutant](https://github.com/abhimutant))
- Static Routes Feature [\#468] (https://github.com/nutanix/terraform-provider-nutanix/pull/468) [abhimutant](https://github.com/abhimutant))
- Examples and testcases related to Networking Features [\#470] (https://github.com/nutanix/terraform-provider-nutanix/pull/470) ([abhimutant](https://github.com/abhimutant))

**Implemented enhancements:**

- Support of External and Overlay Subnets [\#460](https://github.com/nutanix/terraform-provider-nutanix/pull/460)


## 1.5.0 (May 7, 2022)

[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.5.0-beta.2...v1.5.0)

**Merged pull requests:**

- Update Acceptance tests workflow to run specific provider tests [\#437](https://github.com/nutanix/terraform-provider-nutanix/pull/437) ([bhatipradeep](https://github.com/bhatipradeep))
- Foundation unit tests and existing test fixes [\#433](https://github.com/nutanix/terraform-provider-nutanix/pull/433) ([bhatipradeep](https://github.com/bhatipradeep))
- Foundation Acceptance tests and minor fixes [\#436](https://github.com/nutanix/terraform-provider-nutanix/pull/436) ([bhatipradeep](https://github.com/bhatipradeep))
- Foundation Central unit tests and acceptance tests [\#439](https://github.com/nutanix/terraform-provider-nutanix/pull/439) ([abhimutant](https://github.com/abhimutant))
- Optimize Image upload to avoid buffering. Add cluster related fields in image upload resource & data source for PC [\#432](https://github.com/nutanix/terraform-provider-nutanix/pull/432) ([bhatipradeep](https://github.com/bhatipradeep))
- fixing karbon docs at registry [\#434](https://github.com/nutanix/terraform-provider-nutanix/pull/434)([abhimutant](https://github.com/abhimutant)) 
- Example for using config downloaded from install.nutanix.com to image nodes[\#444](https://github.com/nutanix/terraform-provider-nutanix/pull/444) ([bhatipradeep](https://github.com/bhatipradeep))
- Add example to pull secrets from hashicorp vault to use them in node imaging [\#431](https://github.com/nutanix/terraform-provider-nutanix/pull/431) ([bhatipradeep](https://github.com/bhatipradeep))


**Fixed Bugs:**

- Image upload failures [\#430](https://github.com/nutanix/terraform-provider-nutanix/issues/430)
- Not able to upload image to a specific cluster [\#428](https://github.com/nutanix/terraform-provider-nutanix/issues/428)
- Doc for Nutanix karbon_cluster_kubeConfig pointing to deadlink in registry [\#429](https://github.com/nutanix/terraform-provider-nutanix/issues/429)


**Closed Issues:**

- Add Example for using config json file downloaded from install.nutanix.com to image nodes using terraform nutanix foundation [\#435](https://github.com/nutanix/terraform-provider-nutanix/issues/435)
- Image upload failures [\#430](https://github.com/nutanix/terraform-provider-nutanix/issues/430)
- Not able to upload image to a specific cluster [\#428](https://github.com/nutanix/terraform-provider-nutanix/issues/428)
- Add IPMI validations for verifying successfull IPMI configuration [\#425](https://github.com/nutanix/terraform-provider-nutanix/issues/425)
- Add examples of vault integration with node imaging using foundation [\#424](https://github.com/nutanix/terraform-provider-nutanix/issues/424)


## 1.5.0-beta.2 (April 12, 2022)

[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.5.0-beta...v1.5.0-beta.2)

**New Feature:**

- Add resources and data sources for Nutanix Foundation Central  [\#422](https://github.com/nutanix/terraform-provider-nutanix/pull/422)

    New Data Sources :
    -   nutanix_foundation_central_api_keys
    -   nutanix_foundation_central_list_api_keys
    -   nutanix_foundation_central_imaged_nodes_list
    -   nutanix_foundation_central_imaged_clusters_list
    -   nutanix_foundation_central_cluster_details
    -   nutanix_foundation_central_imaged_node_details

    
    New Resources :
    -   nutanix_foundation_central_image_cluster
    -   nutanix_foundation_central_api_keys

    New Modules : 
    -   aos-based-node-imaging/node-serials-filter
    -   manual-mode-imaging


## 1.5.0-beta (April 1, 2022)

[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.4.1...v1.5.0-beta)

**New Feature:**

- Add resources and data sources for Nutanix Foundation  [\#413](https://github.com/nutanix/terraform-provider-nutanix/issues/413)

    New Data Sources :
    -   nutanix_foundation_nos_packages
    -   nutanix_foundation_hypervisor_isos
    -   nutanix_foundation_discover_nodes
    -   nutanix_foundation_node_network_details
    
    New Resources :
    -   nutanix_foundation_image_nodes
    -   nutanix_foundation_ipmi_config
    -   nutanix_foundation_image

    New Modules : 
    -   aos-based-node-imaging/node-serials-filter
    -   discover-nodes-network-details/node-serials-filter
    -   manual-mode-imaging

**Merged pull requests:**

- Feature: Foundation datasources, resources, modules & documentation [\#410](https://github.com/nutanix/terraform-provider-nutanix/pull/410) ([bhatipradeep](https://github.com/bhatipradeep)) ([abhimutant](https://github.com/abhimutant))

> For the 1.5.0-beta release of the provider it will have N-1 compatibility with the Nutanix Foundation. This release was tested against Foundation versions v5.2 and v5.1.1

## 1.4.1 (March 24, 2022)

[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.4.0...v1.4.1)


**Fixed Bugs**

- Network Security Rules: ip subnet is passed as an empty object in the API to create network security which causes API to get failed [\#399](https://github.com/nutanix/terraform-provider-nutanix/issues/399)


**Merged pull requests:**

- Network Security Rules: ip subnet is passed as an empty object in the API to create network security which causes API to get failed [\#348](https://github.com/nutanix/terraform-provider-nutanix/pull/397)([siddharth-kulshrestha](https://github.com/siddharth-kulshrestha))



## 1.4.0 (March 17, 2022)

[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.3.0...v1.4.0)

**Implemented enhancements:**

- Datasources for service groups along with test cases [\#353](https://github.com/nutanix/terraform-provider-nutanix/issues/353)
- Set default value of disk_list.device_properties.device_type to DISK [\#157](https://github.com/nutanix/terraform-provider-nutanix/pull/157)
- Add num_queues parameter for virtual machine NICs [\#324](https://github.com/nutanix/terraform-provider-nutanix/pull/324)
- Add support for linking service and address groups to network security rules [\#345](https://github.com/nutanix/terraform-provider-nutanix/pull/345)


**Fixed Bugs**

- Karbon 2.4 Cluster creation failing [\#373](https://github.com/nutanix/terraform-provider-nutanix/issues/373)
- Use correct API endpoint for category queries [\#348](https://github.com/nutanix/terraform-provider-nutanix/pull/348)
- Terraform crashed while creating VM in Nutanix [\#205](https://github.com/nutanix/terraform-provider-nutanix/issues/205)
- nutanix_recovery_plan "Duplicate network mapping" [\#386](https://github.com/nutanix/terraform-provider-nutanix/issues/386)

**Closed issues:**

- fixed unchecked gets in getMetadataAttributes [\#206](https://github.com/nutanix/terraform-provider-nutanix/pull/206)
- Support defining categories_mapping for a VM [\#321](https://github.com/nutanix/terraform-provider-nutanix/issues/321)
- Add a complete example of parameter nutanix_guest_tools [\#299](https://github.com/nutanix/terraform-provider-nutanix/issues/299)

**Merged pull requests:**

- Use correct API endpoint for category queries [\#348](https://github.com/nutanix/terraform-provider-nutanix/pull/348)([embik](https://github.com/embik))
- Set default value of disk_list.device_properties.device_type to DISK [\#157](https://github.com/nutanix/terraform-provider-nutanix/pull/157) ([maxaudron](https://github.com/maxaudron))
- Added singular and plural datasources for service groups with test cases [\#363](https://github.com/nutanix/terraform-provider-nutanix/pull/363) ([siddharth-kulshrestha](https://github.com/siddharth-kulshrestha))
- Add acceptance tests [\#355](https://github.com/nutanix/terraform-provider-nutanix/pull/355) ([bhatipradeep](https://github.com/bhatipradeep))
- Change default node pool names for karbon cluster creation [\#381](https://github.com/nutanix/terraform-provider-nutanix/pull/381) [bhatipradeep](https://github.com/bhatipradeep))
- bug fix for issue-205 [\#380](https://github.com/nutanix/terraform-provider-nutanix/pull/380) ([abhimutant](https://github.com/abhimutant))
- Task/m 324 check [\#385](https://github.com/nutanix/terraform-provider-nutanix/pull/385) ([abhimutant](https://github.com/abhimutant))
- Security rules with service group and address groups [\#384](https://github.com/nutanix/terraform-provider-nutanix/pull/384) ([siddharth-kulshrestha](https://github.com/siddharth-kulshrestha))
- Fix bug where networkmappings are incorrectly referenced in recovery plan [\#392](https://github.com/nutanix/terraform-provider-nutanix/pull/392) ([yannickstruyf3](https://github.com/yannickstruyf3))



## 1.3.0 (Feb 16, 2022)

[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.2.2...v1.3.0)

**Implemented enhancements:**

- Implement Timeout blocks on resource level [\#254](https://github.com/nutanix/terraform-provider-nutanix/issues/254)
- SDK V2 Upgrade: Upgrade Terraform provider to the latest SDK [\#287](https://github.com/nutanix/terraform-provider-nutanix/issues/287)
- Allow Subnet Datasources to filter based on PE ID [\#308](https://github.com/nutanix/terraform-provider-nutanix/issues/308)
- Implement CI/CD pipeline for this project [\#314](https://github.com/nutanix/terraform-provider-nutanix/issues/314)
- Clean up existing tests with dynamic config  [\#315](https://github.com/nutanix/terraform-provider-nutanix/issues/315)
- Develop integration tests with code coverage. [\#316](https://github.com/nutanix/terraform-provider-nutanix/issues/316)


**Fixed Bugs**

- Provider plugin crashes when nutanix_subnet datasource queried with overlay subnet [\#327](https://github.com/nutanix/terraform-provider-nutanix/issues/327)
- `nutanix_projects` and `nutanix_project` throws error when there is external network associated to a project [\#337](https://github.com/nutanix/terraform-provider-nutanix/issues/337)


**Closed issues:**

- Implement Timeout blocks on resource level [\#254](https://github.com/nutanix/terraform-provider-nutanix/issues/254)
- Upgrade Terraform provider to the latest SDK [\#287](https://github.com/nutanix/terraform-provider-nutanix/issues/287)
- Allow Subnet Datasources to filter based on PE ID [\#308](https://github.com/nutanix/terraform-provider-nutanix/issues/308)
- Implement CI/CD pipeline for this project [\#314](https://github.com/nutanix/terraform-provider-nutanix/issues/314)
- Clean up existing tests with dynamic config  [\#315](https://github.com/nutanix/terraform-provider-nutanix/issues/315)
- Develop integration tests with code coverage. [\#316](https://github.com/nutanix/terraform-provider-nutanix/issues/316)
- Provider plugin crashes when nutanix_subnet datasource queried with overlay subnet [\#327](https://github.com/nutanix/terraform-provider-nutanix/issues/327)
- `nutanix_projects` and `nutanix_project` throws error when there is external network associated to a project [\#337](https://github.com/nutanix/terraform-provider-nutanix/issues/337)



**Merged pull requests:**

- Workflow for automated acceptance test cases [\#325](https://github.com/nutanix/terraform-provider-nutanix/pull/325) ([siddharth-kulshrestha](https://github.com/siddharth-kulshrestha))
- update go release to 1.17 [\#317](https://github.com/nutanix/terraform-provider-nutanix/pull/317) ([tuxtof](https://github.com/tuxtof))
- Workflow for automated acceptance test cases [\#325](https://github.com/nutanix/terraform-provider-nutanix/pull/325) ([siddharth-kulshrestha](https://github.com/siddharth-kulshrestha))
- Fix nutanix_project and nutanix_projects datasource to not use "kind" attribute in "external_network_list" [\#335](https://github.com/nutanix/terraform-provider-nutanix/pull/335) ([bhati-pradeep](https://github.com/bhati-pradeep))
- Add check for cluster_reference before set to avoid it in case of overlay subnets in datasource nutanix_subnet [\#328](https://github.com/nutanix/terraform-provider-nutanix/pull/328) ([bhati-pradeep](https://github.com/bhati-pradeep))
- fixed a typo in subnet.html.markdown [\#273](https://github.com/nutanix/terraform-provider-nutanix/pull/273) ([gowatana](https://github.com/gowatana))
- Update subnets.html.markdown [\#293](https://github.com/nutanix/terraform-provider-nutanix/pull/293) ([jastrom85](https://github.com/jastrom85))
- Add filter by cluster uuid in subnet datasource [\#323](https://github.com/nutanix/terraform-provider-nutanix/pull/323) ([shreevari](https://github.com/shreevari))
- Modify tests and example to use is_vcpu_hard_pinned for nutanix_virtual_machine [\#338](https://github.com/nutanix/terraform-provider-nutanix/pull/338) ([bhati-pradeep](https://github.com/bhati-pradeep))
- Add resources for address groups and service groups [\#322](https://github.com/nutanix/terraform-provider-nutanix/pull/322) ([svalabs](https://github.com/svalabs))
- Service group fix [\#340](https://github.com/nutanix/terraform-provider-nutanix/pull/340) ([abhimutant](https://github.com/abhimutant))
- V2 SDK upgrade [\#332](https://github.com/nutanix/terraform-provider-nutanix/pull/332) ([siddharth-kulshrestha](https://github.com/siddharth-kulshrestha))
- Add vCPU hard pinning  [\#307](https://github.com/nutanix/terraform-provider-nutanix/pull/307) ([basraayman](https://github.com/basraayman))
- bugfix 254 - timeout on resource level [\#333](https://github.com/nutanix/terraform-provider-nutanix/pull/333) ([venkatavivek-ntnx](https://github.com/venkatavivek-ntnx))
- Fix existing examples. Add new examples [\#346](https://github.com/nutanix/terraform-provider-nutanix/pull/346) ([bhati-pradeep](https://github.com/bhati-pradeep))
- Automation for manual testcases [\#334](https://github.com/nutanix/terraform-provider-nutanix/issues/334) ([abhimutant](https://github.com/abhimutant))



## 1.2.2 (Nov 29, 2021)

[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.2.1...v1.2.2)

**Closed issues:**

- Cloud_init and sysprep CDROMs being detected and destroyed forcing VM reboot. [\#285](https://github.com/nutanix/terraform-provider-nutanix/issues/285)
- Provider crash when using NGT credentials attribute in nutanix_virtual_machine resource type. [\#297](https://github.com/nutanix/terraform-provider-nutanix/issues/297)

**Merged pull requests:**

- Add support for CPU passthrough [\#288](https://github.com/nutanix/terraform-provider-nutanix/pull/288) ([svalabs](https://github.com/svalabs))



## 1.2.1 (Feb 01, 2021)

[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.1.1...HEAD)

**Closed issues:**

- Terraform crashed while creating VM in Nutanix [\#205](https://github.com/nutanix/terraform-provider-nutanix/issues/205)
- nutanix provider doesn't appear in terraform registry [\#166](https://github.com/nutanix/terraform-provider-nutanix/issues/166)
- Option not available to change the BIOS type [\#163](https://github.com/nutanix/terraform-provider-nutanix/issues/163)
- Need a user/group data source to get uuid [\#142](https://github.com/nutanix/terraform-provider-nutanix/issues/142)
- Missing user role association with nutanix\_project ressource [\#141](https://github.com/nutanix/terraform-provider-nutanix/issues/141)
- Support for Secure VDI Groups \(VDI Policy\) in nutanix\_network\_security\_rule [\#132](https://github.com/nutanix/terraform-provider-nutanix/issues/132)
- Provider does not recover from error [\#131](https://github.com/nutanix/terraform-provider-nutanix/issues/131)
- During a resource adjustment the VM is not shut down using ACPI [\#96](https://github.com/nutanix/terraform-provider-nutanix/issues/96)
- DR Runbook support [\#22](https://github.com/nutanix/terraform-provider-nutanix/issues/22)
- Nutanix Karbon interaction [\#6](https://github.com/nutanix/terraform-provider-nutanix/issues/6)

**Merged pull requests:**

- Add note about Prism version api we tested for v1.2.0 release [\#246](https://github.com/nutanix/terraform-provider-nutanix/pull/246) ([marinsalinas](https://github.com/marinsalinas))
- Update README and Website for  v1.2.0 release. [\#245](https://github.com/nutanix/terraform-provider-nutanix/pull/245) ([marinsalinas](https://github.com/marinsalinas))
- persist uuid after creation [\#244](https://github.com/nutanix/terraform-provider-nutanix/pull/244) ([yannickstruyf3](https://github.com/yannickstruyf3))
- V1.2.0 beta [\#243](https://github.com/nutanix/terraform-provider-nutanix/pull/243) ([marinsalinas](https://github.com/marinsalinas))
- test: changed data for variables of azurl for protection rules [\#242](https://github.com/nutanix/terraform-provider-nutanix/pull/242) ([coderGo93](https://github.com/coderGo93))
- V1.2.0 fix testacc [\#241](https://github.com/nutanix/terraform-provider-nutanix/pull/241) ([marinsalinas](https://github.com/marinsalinas))
- fixed bug where datasource project search by name was empty [\#240](https://github.com/nutanix/terraform-provider-nutanix/pull/240) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Bugfix/remove ide3 dependency tests [\#239](https://github.com/nutanix/terraform-provider-nutanix/pull/239) ([marinsalinas](https://github.com/marinsalinas))
- Added missing information for resource and datasource of project [\#238](https://github.com/nutanix/terraform-provider-nutanix/pull/238) ([coderGo93](https://github.com/coderGo93))
- Added filter by name in datasource of project [\#237](https://github.com/nutanix/terraform-provider-nutanix/pull/237) ([coderGo93](https://github.com/coderGo93))
- Added filter by name in datasource of access control policy [\#236](https://github.com/nutanix/terraform-provider-nutanix/pull/236) ([coderGo93](https://github.com/coderGo93))
- Task completion project [\#234](https://github.com/nutanix/terraform-provider-nutanix/pull/234) ([coderGo93](https://github.com/coderGo93))
- adding a VM to a project does not require a reboot [\#232](https://github.com/nutanix/terraform-provider-nutanix/pull/232) ([yannickstruyf3](https://github.com/yannickstruyf3))
- removed project\_ref [\#231](https://github.com/nutanix/terraform-provider-nutanix/pull/231) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Bugfix/v1.2.0 beta doc review yst [\#229](https://github.com/nutanix/terraform-provider-nutanix/pull/229) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Karbon Base [\#228](https://github.com/nutanix/terraform-provider-nutanix/pull/228) ([marinsalinas](https://github.com/marinsalinas))
- Bugfix/remove ide3 dependency [\#227](https://github.com/nutanix/terraform-provider-nutanix/pull/227) ([yannickstruyf3](https://github.com/yannickstruyf3))
- added machine\_type to data source [\#226](https://github.com/nutanix/terraform-provider-nutanix/pull/226) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Add testacc for karbon resources and data sources [\#222](https://github.com/nutanix/terraform-provider-nutanix/pull/222) ([marinsalinas](https://github.com/marinsalinas))
- Karbon docs [\#221](https://github.com/nutanix/terraform-provider-nutanix/pull/221) ([coderGo93](https://github.com/coderGo93))
- VDI Support [\#220](https://github.com/nutanix/terraform-provider-nutanix/pull/220) ([marinsalinas](https://github.com/marinsalinas))
- added permissions for roles [\#219](https://github.com/nutanix/terraform-provider-nutanix/pull/219) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Feat: add bios\_type support to VM resource and datasource [\#217](https://github.com/nutanix/terraform-provider-nutanix/pull/217) ([marinsalinas](https://github.com/marinsalinas))
- Protection rule and Recovery plan\(DR Runbook\) resources and datasources [\#216](https://github.com/nutanix/terraform-provider-nutanix/pull/216) ([coderGo93](https://github.com/coderGo93))
- fixed bug that occured when updating the permissions of a role [\#215](https://github.com/nutanix/terraform-provider-nutanix/pull/215) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Bugfix/role update [\#214](https://github.com/nutanix/terraform-provider-nutanix/pull/214) ([yannickstruyf3](https://github.com/yannickstruyf3))
- renamed ds attributes and prints [\#213](https://github.com/nutanix/terraform-provider-nutanix/pull/213) ([yannickstruyf3](https://github.com/yannickstruyf3))
- changed print messages and added filter based on DN [\#212](https://github.com/nutanix/terraform-provider-nutanix/pull/212) ([yannickstruyf3](https://github.com/yannickstruyf3))
- changed if statement in ds role and modified the conflictswith [\#195](https://github.com/nutanix/terraform-provider-nutanix/pull/195) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Roles [\#181](https://github.com/nutanix/terraform-provider-nutanix/pull/181) ([coderGo93](https://github.com/coderGo93))
- User Resource and Data Sources. [\#179](https://github.com/nutanix/terraform-provider-nutanix/pull/179) ([marinsalinas](https://github.com/marinsalinas))
- Access control policies [\#175](https://github.com/nutanix/terraform-provider-nutanix/pull/175) ([coderGo93](https://github.com/coderGo93))
## [1.1.1](https://github.com/nutanix/terraform-provider-nutanix/tree/v1.1.1) (2020-11-30)
[Full Changelog](https://github.com/nutanix/terraform-provider-nutanix/compare/v1.1.0...v1.1.1)

**Fixed bugs:**

- local resource nutanix\_image [\#182](https://github.com/nutanix/terraform-provider-nutanix/issues/182)
- Terraform crashes when trying to specify boot\_device\_order\_list for a vm resource [\#28](https://github.com/nutanix/terraform-provider-nutanix/issues/28)

**Closed issues:**

- Problem on json unmarshalling to go struct \(MessageResource.message\_list.details of type map\[string\]interface {}\) [\#204](https://github.com/nutanix/terraform-provider-nutanix/issues/204)
- Fix resource Read inconsistencies [\#201](https://github.com/nutanix/terraform-provider-nutanix/issues/201)
- local resource nutanix\_image [\#182](https://github.com/nutanix/terraform-provider-nutanix/issues/182)
- vss\_snapshot\_capable attribute doesn't work [\#171](https://github.com/nutanix/terraform-provider-nutanix/issues/171)
- 'invalid memory address or nil pointer dereference' while getting a VM. [\#168](https://github.com/nutanix/terraform-provider-nutanix/issues/168)
- FYI: Nutanix API Authentication fails when making many requests at once [\#93](https://github.com/nutanix/terraform-provider-nutanix/issues/93)
- 1122 API requests only to list VMs [\#87](https://github.com/nutanix/terraform-provider-nutanix/issues/87)
- Failed to perform NGT ISO mount operation with error code: kNoFreeCdRomSlot [\#83](https://github.com/nutanix/terraform-provider-nutanix/issues/83)
- ip\_address not available as nutanix\_virtual\_machine attribute. [\#19](https://github.com/nutanix/terraform-provider-nutanix/issues/19)
- Unable to obtain Private\_IP as a Output  [\#17](https://github.com/nutanix/terraform-provider-nutanix/issues/17)

**Merged pull requests:**

- chore: fix goreleaser for v1.1.1 [\#207](https://github.com/nutanix/terraform-provider-nutanix/pull/207) ([marinsalinas](https://github.com/marinsalinas))
- added return nil in read function [\#203](https://github.com/nutanix/terraform-provider-nutanix/pull/203) ([yannickstruyf3](https://github.com/yannickstruyf3))
- fixed source\_path image\_type issue [\#202](https://github.com/nutanix/terraform-provider-nutanix/pull/202) ([yannickstruyf3](https://github.com/yannickstruyf3))
- fixed wrong status\_code check in client \(404 should be 401\) [\#200](https://github.com/nutanix/terraform-provider-nutanix/pull/200) ([yannickstruyf3](https://github.com/yannickstruyf3))
- GitHub actions [\#198](https://github.com/nutanix/terraform-provider-nutanix/pull/198) ([marinsalinas](https://github.com/marinsalinas))
- chore: add note to network\_security\_rule [\#196](https://github.com/nutanix/terraform-provider-nutanix/pull/196) ([marinsalinas](https://github.com/marinsalinas))
- chore: fix linting issues [\#194](https://github.com/nutanix/terraform-provider-nutanix/pull/194) ([marinsalinas](https://github.com/marinsalinas))
- Fix Testacc [\#193](https://github.com/nutanix/terraform-provider-nutanix/pull/193) ([marinsalinas](https://github.com/marinsalinas))
- fix: missing variable initialization [\#192](https://github.com/nutanix/terraform-provider-nutanix/pull/192) ([marinsalinas](https://github.com/marinsalinas))
- Migrate to Terraform Plugin SDK v1 [\#176](https://github.com/nutanix/terraform-provider-nutanix/pull/176) ([marinsalinas](https://github.com/marinsalinas))
- Fix issue \#168, unchecked nil assignment [\#169](https://github.com/nutanix/terraform-provider-nutanix/pull/169) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Added name argument for datasource of cluster [\#165](https://github.com/nutanix/terraform-provider-nutanix/pull/165) ([coderGo93](https://github.com/coderGo93))
- Removed volume\_group documentations [\#160](https://github.com/nutanix/terraform-provider-nutanix/pull/160) ([PacoDw](https://github.com/PacoDw))
- Updated Readme file and changelog [\#154](https://github.com/nutanix/terraform-provider-nutanix/pull/154) ([PacoDw](https://github.com/PacoDw))
- provide better error reporting when invalid nutanix credentials are used [\#148](https://github.com/nutanix/terraform-provider-nutanix/pull/148) ([scott-grimes](https://github.com/scott-grimes))
- Initialize BootConfig struct before the child attributes [\#53](https://github.com/nutanix/terraform-provider-nutanix/pull/53) ([marinsalinas](https://github.com/marinsalinas))



## 1.1.0 (July 02, 2020)

**Implemented enhancements:**

- Boot device order list was limitated to set 1 item until fix issue #28, documentation was updated about it [\#152](https://github.com/terraform-providers/terraform-provider-nutanix/pull/152) ([PacoDw](https://github.com/PacoDw))
- Added Host datasources [\#116](https://github.com/terraform-providers/terraform-provider-nutanix/pull/116) ([PacoDw](https://github.com/PacoDw))
- Added datasource of subnets [\#115](https://github.com/terraform-providers/terraform-provider-nutanix/pull/115) ([coderGo93](https://github.com/coderGo93))
- Validated the terraform configuration adding new test case to validate it [\#114](https://github.com/terraform-providers/terraform-provider-nutanix/pull/114) ([PacoDw](https://github.com/PacoDw))
- Implemented new linter version and fixed new version linter issues [\#101](https://github.com/terraform-providers/terraform-provider-nutanix/pull/101) ([marinsalinas](https://github.com/marinsalinas))
- Updated docs by adding should_force_translated for metadata for every datasource/resource available [\#140](https://github.com/terraform-providers/terraform-provider-nutanix/pull/140) ([coderGo93](https://github.com/coderGo93))
- Documentations for host data sources were added [\#139](https://github.com/terraform-providers/terraform-provider-nutanix/pull/139) ([coderGo93](https://github.com/coderGo93))
- The storage_conntainer was added as a new attribute into the disk_list to reference a container in a VM [\#138](https://github.com/terraform-providers/terraform-provider-nutanix/pull/138) ([PacoDw](https://github.com/PacoDw))
- Added more information about Sysprep for install type in the documentation [\#133](https://github.com/terraform-providers/terraform-provider-nutanix/pull/133) ([coderGo93](https://github.com/coderGo93))
- Added preCheck function to verify that env variables were set [\#103](https://github.com/terraform-providers/terraform-provider-nutanix/pull/103) ([PacoDw](https://github.com/PacoDw))
- Feature/cluster datasource name [\#100](https://github.com/terraform-providers/terraform-provider-nutanix/pull/100) ([yannickstruyf3](https://github.com/yannickstruyf3))

**Fixed bugs:**

- Unable to create vm using guest\_customization\_cloud\_init\_custom\_key\_values [\#58](https://github.com/terraform-providers/terraform-provider-nutanix/issues/58)
- Fixed the behavior of the use\_hot\_add attribute to reboot the VM correctly [\#149](https://github.com/terraform-providers/terraform-provider-nutanix/pull/149) ([PacoDw](https://github.com/PacoDw))
- Fixed storage container attribute changing it to computed [\#147](https://github.com/terraform-providers/terraform-provider-nutanix/pull/147) ([PacoDw](https://github.com/PacoDw))
- Manually deleting VM and running apply results in entity not found. vm exists code removed [\#144](https://github.com/terraform-providers/terraform-provider-nutanix/pull/144) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Removed check to make sure data_source_reference and disk_size_bytes are mutually exclusive [\#137](https://github.com/terraform-providers/terraform-provider-nutanix/pull/137) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Added new parameter for VM use\_hot\_add [\#136](https://github.com/terraform-providers/terraform-provider-nutanix/pull/136) ([coderGo93](https://github.com/coderGo93))
- Improve error handling on incorrect API calls [\#134](https://github.com/terraform-providers/terraform-provider-nutanix/pull/134) ([PacoDw](https://github.com/PacoDw))
- The filter for subnets was incorrect. Filtering on name is not via subnet_name but name [\#129](https://github.com/terraform-providers/terraform-provider-nutanix/pull/129) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Fixed Security rules idempotency error: Error was thrown when security rule was delete via Prism and running a new Terraform run afterwards [\#127](https://github.com/terraform-providers/terraform-provider-nutanix/pull/127) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Added validation to avoid nil pointer error [\#122](https://github.com/terraform-providers/terraform-provider-nutanix/pull/122) ([PacoDw](https://github.com/PacoDw))
- Fixed timeout issue and categories reboot [\#120](https://github.com/terraform-providers/terraform-provider-nutanix/pull/120) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Bugfix/cloudinit final [\#111](https://github.com/terraform-providers/terraform-provider-nutanix/pull/111) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Fixed guest_customization_cloud_init_custom_key_values to create the VM correctly [\#102](https://github.com/terraform-providers/terraform-provider-nutanix/pull/102) ([yannickstruyf3](https://github.com/yannickstruyf3))

**Closed issues:**

- ENTITY\_NOT\_FOUND Error on apply after manual VM deletion [\#143](https://github.com/terraform-providers/terraform-provider-nutanix/issues/143)
- Provider prevents disk resizing at VM creation time [\#130](https://github.com/terraform-providers/terraform-provider-nutanix/issues/130)
- Datasource subnets subnet\_name throwing error [\#128](https://github.com/terraform-providers/terraform-provider-nutanix/issues/128)
- data source nutanix\_network\_security\_rule expected type issue [\#126](https://github.com/terraform-providers/terraform-provider-nutanix/issues/126)
- Categories assignation issue [\#119](https://github.com/terraform-providers/terraform-provider-nutanix/issues/119)
- The nutanix\_subnet SIGSEGV on network rename [\#118](https://github.com/terraform-providers/terraform-provider-nutanix/issues/118)
- Improve error handling on incorrect API calls [\#112](https://github.com/terraform-providers/terraform-provider-nutanix/issues/112)
- The nutanix\_virtual\_machine resource won't allow setting power\_state [\#98](https://github.com/terraform-providers/terraform-provider-nutanix/issues/98)
- Unable to modify a VM with learned IP address \(DHCP\) fails with "IP address with type 'LEARNED' not allowed" [\#97](https://github.com/terraform-providers/terraform-provider-nutanix/issues/97)
- How do I use an unattend.xml file? [\#95](https://github.com/terraform-providers/terraform-provider-nutanix/issues/95)
- power\_state trigger a syntax error. [\#94](https://github.com/terraform-providers/terraform-provider-nutanix/issues/94)
- Undefined property is\_connected of class vm\_nic detected [\#90](https://github.com/terraform-providers/terraform-provider-nutanix/issues/90)
- Update Documentation for guest\_customization\_sysprep Attribute [\#89](https://github.com/terraform-providers/terraform-provider-nutanix/issues/89)
- CPU and RAM change does not restart VM automatically [\#86](https://github.com/terraform-providers/terraform-provider-nutanix/issues/86)
- The user\_data for Windows VM [\#84](https://github.com/terraform-providers/terraform-provider-nutanix/issues/84)
- Cloning from the Image Server leads to a change after a second terraform apply [\#82](https://github.com/terraform-providers/terraform-provider-nutanix/issues/82)
- Hotplugging CPU and RAM [\#79](https://github.com/terraform-providers/terraform-provider-nutanix/issues/79)
- Cannot specify target container when adding disks to a virtual machine resource [\#78](https://github.com/terraform-providers/terraform-provider-nutanix/issues/78)
- Cannot use nutanix\_subnets data source [\#73](https://github.com/terraform-providers/terraform-provider-nutanix/issues/73)
- Add datasource  type "hosts" for api/nutanix/v3/hosts endpoint [\#70](https://github.com/terraform-providers/terraform-provider-nutanix/issues/70)
- Using user\_data on resource.nutanix\_virtual\_machine yields immediate diff after initial apply [\#69](https://github.com/terraform-providers/terraform-provider-nutanix/issues/69)
- Cannot list ip addresses of when creating multiple virtual machine resources  [\#63](https://github.com/terraform-providers/terraform-provider-nutanix/issues/63)
- Provider needs to support data\_source\_reference attribute for nutanix\_image resource [\#52](https://github.com/terraform-providers/terraform-provider-nutanix/issues/52)
- Terraform Unable to use Provisioners in VM Resource when DHCP is used for IP Address [\#49](https://github.com/terraform-providers/terraform-provider-nutanix/issues/49)
- Cloning from a VM on AHV [\#35](https://github.com/terraform-providers/terraform-provider-nutanix/issues/35)

**Merged pull requests:**

- V1.1.0 rc1 [\#150](https://github.com/terraform-providers/terraform-provider-nutanix/pull/150) ([PacoDw](https://github.com/PacoDw))
- Disabled project resource, tests fixed and issue fixed [\#146](https://github.com/terraform-providers/terraform-provider-nutanix/pull/146) ([PacoDw](https://github.com/PacoDw))
- Improved travis [\#117](https://github.com/terraform-providers/terraform-provider-nutanix/pull/117) ([PacoDw](https://github.com/PacoDw))
- Bugfix device\_properties in a disk\_list [\#106](https://github.com/terraform-providers/terraform-provider-nutanix/pull/106) ([yannickstruyf3](https://github.com/yannickstruyf3))
- Reduce the amount of API requests and improvements [\#92](https://github.com/terraform-providers/terraform-provider-nutanix/pull/92) ([maxaudron](https://github.com/maxaudron))
- Implement session based authentification [\#88](https://github.com/terraform-providers/terraform-provider-nutanix/pull/88) ([maxaudron](https://github.com/maxaudron))
- imrpoved wesite removing # [\#85](https://github.com/terraform-providers/terraform-provider-nutanix/pull/85) ([mritzmann](https://github.com/mritzmann))
- Add support for mapstructure decode tag for reusability [\#81](https://github.com/terraform-providers/terraform-provider-nutanix/pull/81) ([JRemitz](https://github.com/JRemitz))
- Add missing API fields for structs [\#80](https://github.com/terraform-providers/terraform-provider-nutanix/pull/80) ([JRemitz](https://github.com/JRemitz))

## 1.0.2 (September 05, 2019)

**Fixed bugs:**

- user\_data typos cause panic, ForceNew for guest\_customization fields [\#67](https://github.com/terraform-providers/terraform-provider-nutanix/issues/67)
- Example config is not valid for power\_state [\#9](https://github.com/terraform-providers/terraform-provider-nutanix/issues/9)

**Closed issues:**

- Do not shutdown machine on certain changes [\#74](https://github.com/terraform-providers/terraform-provider-nutanix/issues/74)
- Update compatibility matrix for TF provider [\#71](https://github.com/terraform-providers/terraform-provider-nutanix/issues/71)
- Terraform 0.12 compatibility [\#66](https://github.com/terraform-providers/terraform-provider-nutanix/issues/66)
- Unable to specify category key names dynamically in resource:nutanix\_virtual\_machine [\#61](https://github.com/terraform-providers/terraform-provider-nutanix/issues/61)
- data source nutanix\_virtual\_machine fails with "Invalid address to set: \[\]string{"nic\_list", "0", "is\_connected"}" [\#57](https://github.com/terraform-providers/terraform-provider-nutanix/issues/57)
- cannot unmarshal string into GO struct field MessageResource.details of the type map \[string\]interface{} [\#44](https://github.com/terraform-providers/terraform-provider-nutanix/issues/44)
- Error when re-applying a plan [\#41](https://github.com/terraform-providers/terraform-provider-nutanix/issues/41)
- unable to spin vm with error "'api\_version' is a required property" [\#36](https://github.com/terraform-providers/terraform-provider-nutanix/issues/36)
- Change VM caused disk being deleted [\#34](https://github.com/terraform-providers/terraform-provider-nutanix/issues/34)
- Unable to change VM resources when it has a network connected [\#33](https://github.com/terraform-providers/terraform-provider-nutanix/issues/33)
- Consider replacing satori/go.uuid [\#31](https://github.com/terraform-providers/terraform-provider-nutanix/issues/31)

**Merged pull requests:**

- chore: update changelog for 1.0.2 version [\#77](https://github.com/terraform-providers/terraform-provider-nutanix/pull/77) ([marinsalinas](https://github.com/marinsalinas))
- Dynamic category name support [\#76](https://github.com/terraform-providers/terraform-provider-nutanix/pull/76) ([marinsalinas](https://github.com/marinsalinas))
- Patch to allow some changes to be hotplug [\#75](https://github.com/terraform-providers/terraform-provider-nutanix/pull/75) ([Jorge-Holgado](https://github.com/Jorge-Holgado))
- Terraform 0.12 Provider Support. [\#72](https://github.com/terraform-providers/terraform-provider-nutanix/pull/72) ([marinsalinas](https://github.com/marinsalinas))
- Fix typos causing panic on userdata change, ForceNew on change [\#68](https://github.com/terraform-providers/terraform-provider-nutanix/pull/68) ([rxacevedo](https://github.com/rxacevedo))
- Refactor: change update workflow to use spec instead status in virtual machine resource. [\#62](https://github.com/terraform-providers/terraform-provider-nutanix/pull/62) ([marinsalinas](https://github.com/marinsalinas))
- Added fix for Issue \#57 [\#60](https://github.com/terraform-providers/terraform-provider-nutanix/pull/60) ([chandru-tkc](https://github.com/chandru-tkc))
- added VMNic.IsConnected to preFillResUpdateRequest [\#59](https://github.com/terraform-providers/terraform-provider-nutanix/pull/59) ([switchboardOp](https://github.com/switchboardOp))


## 1.0.1 (May 01, 2019)

**Implemented enhancements:**

- Ability to resize a vdisk when deploying from image resource [\#23](https://github.com/terraform-providers/terraform-provider-nutanix/issues/23)
- nutanix\_image: long running image create completes successfully at 10 minutes, without getting success call from API [\#20](https://github.com/terraform-providers/terraform-provider-nutanix/issues/20)
- Datasources for categories [\#8](https://github.com/terraform-providers/terraform-provider-nutanix/issues/8)

**Fixed bugs:**

- nutanix\_image: long running image create completes successfully at 10 minutes, without getting success call from API [\#20](https://github.com/terraform-providers/terraform-provider-nutanix/issues/20)
- nutanix\_virtual\_machine, ip\_endpoint\_list not correctly working [\#12](https://github.com/terraform-providers/terraform-provider-nutanix/issues/12)
- nutanix\_clusters data source makes Terraform crash [\#10](https://github.com/terraform-providers/terraform-provider-nutanix/issues/10)

**Closed issues:**

- Terraform error trying to use nutanix\_virtual\_machine data source [\#30](https://github.com/terraform-providers/terraform-provider-nutanix/issues/30)
- \[PROPOSAL\] Switch to Go Modules [\#29](https://github.com/terraform-providers/terraform-provider-nutanix/issues/29)
- hard code \(kind\) to be appropriate kind for each resource [\#27](https://github.com/terraform-providers/terraform-provider-nutanix/issues/27)
- Implement HTTP proxy capability at provider level [\#26](https://github.com/terraform-providers/terraform-provider-nutanix/issues/26)
- Ability to add a serial port to a VM [\#25](https://github.com/terraform-providers/terraform-provider-nutanix/issues/25)
- Add additional acceptance test configurations for nutanix resources and data sources [\#24](https://github.com/terraform-providers/terraform-provider-nutanix/issues/24)
- Example main.tf - guest\_customization\_cloud\_init does not work [\#21](https://github.com/terraform-providers/terraform-provider-nutanix/issues/21)
- Can't apply just a single resource "nutanix\_image" [\#18](https://github.com/terraform-providers/terraform-provider-nutanix/issues/18)

**Merged pull requests:**

- update changelog for v1.0.1 [\#56](https://github.com/terraform-providers/terraform-provider-nutanix/pull/56) ([marinsalinas](https://github.com/marinsalinas))
- fix category key data source basic test [\#55](https://github.com/terraform-providers/terraform-provider-nutanix/pull/55) ([marinsalinas](https://github.com/marinsalinas))
- Vdisk update \#23 [\#54](https://github.com/terraform-providers/terraform-provider-nutanix/pull/54) ([marinsalinas](https://github.com/marinsalinas))
- Data Source for categories [\#51](https://github.com/terraform-providers/terraform-provider-nutanix/pull/51) ([marinsalinas](https://github.com/marinsalinas))
- chore: update cibuild make formula to build for any OS \(windows, linu… [\#50](https://github.com/terraform-providers/terraform-provider-nutanix/pull/50) ([marinsalinas](https://github.com/marinsalinas))
- chore: fix bad urls in readme.md file [\#48](https://github.com/terraform-providers/terraform-provider-nutanix/pull/48) ([marinsalinas](https://github.com/marinsalinas))
- chore: fix bad urls in readme.md file [\#47](https://github.com/terraform-providers/terraform-provider-nutanix/pull/47) ([marinsalinas](https://github.com/marinsalinas))
- Add HTTP Proxy capability [\#46](https://github.com/terraform-providers/terraform-provider-nutanix/pull/46) ([marinsalinas](https://github.com/marinsalinas))
- Add Serial Port support in VM resource and data source [\#45](https://github.com/terraform-providers/terraform-provider-nutanix/pull/45) ([marinsalinas](https://github.com/marinsalinas))
- Fix: nutanix\_guest\_tools attributes [\#43](https://github.com/terraform-providers/terraform-provider-nutanix/pull/43) ([marinsalinas](https://github.com/marinsalinas))
- Image wait timeout \#20 [\#42](https://github.com/terraform-providers/terraform-provider-nutanix/pull/42) ([marinsalinas](https://github.com/marinsalinas))
- refactor: add nic\_list\_status to separate the computed values [\#40](https://github.com/terraform-providers/terraform-provider-nutanix/pull/40) ([marinsalinas](https://github.com/marinsalinas))
- Remove kind as a argument in cluster\_reference and metadata. [\#39](https://github.com/terraform-providers/terraform-provider-nutanix/pull/39) ([marinsalinas](https://github.com/marinsalinas))
- Fix clusters Data Source makes tf crash \#10 [\#38](https://github.com/terraform-providers/terraform-provider-nutanix/pull/38) ([marinsalinas](https://github.com/marinsalinas))
- \[MODULES\] Switch to Go Modules [\#37](https://github.com/terraform-providers/terraform-provider-nutanix/pull/37) ([appilon](https://github.com/appilon))
- website: change guest\_customization\_cloud\_init argument reference [\#32](https://github.com/terraform-providers/terraform-provider-nutanix/pull/32) ([marinsalinas](https://github.com/marinsalinas))
- \[AUTOMATED\] Upgrade to Go 1.11 [\#16](https://github.com/terraform-providers/terraform-provider-nutanix/pull/16) ([appilon](https://github.com/appilon))
- Update docs to reflect removal of network security rule resource [\#14](https://github.com/terraform-providers/terraform-provider-nutanix/pull/14) ([alias-dev](https://github.com/alias-dev))
- correcting example [\#7](https://github.com/terraform-providers/terraform-provider-nutanix/pull/7) ([olljanat](https://github.com/olljanat))
- Fix Spell and style [\#3](https://github.com/terraform-providers/terraform-provider-nutanix/pull/3) ([ryujisnote](https://github.com/ryujisnote))
- fix cluster datasource\(s\) if no http/smtp credentials are configured [\#1](https://github.com/terraform-providers/terraform-provider-nutanix/pull/1) ([simonfuhrer](https://github.com/simonfuhrer))

## 1.0.0 (October 09, 2018)

**Implemented enhancements:**

- Configure Test Coverage via Code Climate [\#112](https://github.com/nutanix/terraform-provider-nutanix/issues/112)
- Add Cluster DS [\#68](https://github.com/nutanix/terraform-provider-nutanix/issues/68)

**Fixed bugs:**

- client/v3/v3\_service.go - clean up bunch of TODO's [\#127](https://github.com/nutanix/terraform-provider-nutanix/issues/127)
- Incorporate gofmt simplify \(gofmt -s\) [\#91](https://github.com/nutanix/terraform-provider-nutanix/issues/91)
- Fix Go Report Card issues with latest merge [\#90](https://github.com/nutanix/terraform-provider-nutanix/issues/90)

**Closed issues:**

- Get test coverage up \(at least over 50%\) [\#141](https://github.com/nutanix/terraform-provider-nutanix/issues/141)
- Crash during VM change [\#132](https://github.com/nutanix/terraform-provider-nutanix/issues/132)
- Terraform crash during vm create [\#119](https://github.com/nutanix/terraform-provider-nutanix/issues/119)
- Feature requests: max length for metadata / userdata in guest customization [\#115](https://github.com/nutanix/terraform-provider-nutanix/issues/115)
- Terraform crash during vm re-create [\#114](https://github.com/nutanix/terraform-provider-nutanix/issues/114)
- Terraform client fails to read cloud init values [\#111](https://github.com/nutanix/terraform-provider-nutanix/issues/111)
- Volume Groups DS [\#100](https://github.com/nutanix/terraform-provider-nutanix/issues/100)
- Volume Group DS [\#99](https://github.com/nutanix/terraform-provider-nutanix/issues/99)
- Volume Group RS [\#98](https://github.com/nutanix/terraform-provider-nutanix/issues/98)
- Add Subnets DS [\#97](https://github.com/nutanix/terraform-provider-nutanix/issues/97)
- Add Images DS [\#96](https://github.com/nutanix/terraform-provider-nutanix/issues/96)
- Configuration for the Terraform Nutanix Documentation [\#88](https://github.com/nutanix/terraform-provider-nutanix/issues/88)
- Fix "GoLint/Naming/PackageNames" issue in nutanix/config.go [\#87](https://github.com/nutanix/terraform-provider-nutanix/issues/87)
- Fix "SC2046" issue in scripts/errcheck.sh [\#86](https://github.com/nutanix/terraform-provider-nutanix/issues/86)
- Fix "SC2046" issue in scripts/gofmtcheck.sh [\#85](https://github.com/nutanix/terraform-provider-nutanix/issues/85)
- Fix "SC2006" issue in scripts/gofmtcheck.sh [\#84](https://github.com/nutanix/terraform-provider-nutanix/issues/84)
- Fix "SC2006" issue in scripts/changelog-links.sh [\#83](https://github.com/nutanix/terraform-provider-nutanix/issues/83)
- Fix "GoLint/Naming/PackageNames" issue in nutanix/config.go [\#82](https://github.com/nutanix/terraform-provider-nutanix/issues/82)
- Test ticket from Code Climate [\#81](https://github.com/nutanix/terraform-provider-nutanix/issues/81)
- Add Network Security Rules DS [\#80](https://github.com/nutanix/terraform-provider-nutanix/issues/80)
- Add Network Security Rule DS [\#75](https://github.com/nutanix/terraform-provider-nutanix/issues/75)
- Terraform crash with VM and cloud init [\#72](https://github.com/nutanix/terraform-provider-nutanix/issues/72)
- Make all test passing [\#71](https://github.com/nutanix/terraform-provider-nutanix/issues/71)
- Add Cluster RS [\#69](https://github.com/nutanix/terraform-provider-nutanix/issues/69)
- Get Image resource test passing [\#67](https://github.com/nutanix/terraform-provider-nutanix/issues/67)
- implement issue template [\#61](https://github.com/nutanix/terraform-provider-nutanix/issues/61)
- Add codecov.io code coverage reporting [\#57](https://github.com/nutanix/terraform-provider-nutanix/issues/57)
- Add VMS Data Source [\#51](https://github.com/nutanix/terraform-provider-nutanix/issues/51)
- Terraform crash when updating image [\#49](https://github.com/nutanix/terraform-provider-nutanix/issues/49)
- Unclear protocol / error message for image upload [\#48](https://github.com/nutanix/terraform-provider-nutanix/issues/48)
- Cleanup all go report card issues [\#47](https://github.com/nutanix/terraform-provider-nutanix/issues/47)

**Merged pull requests:**

- merge develop june 20 [\#149](https://github.com/nutanix/terraform-provider-nutanix/pull/149) ([JonKohler](https://github.com/JonKohler))
- Add test cases for client and V3  \#141 [\#148](https://github.com/nutanix/terraform-provider-nutanix/pull/148) ([marinsalinas](https://github.com/marinsalinas))
- develop to master nightly [\#147](https://github.com/nutanix/terraform-provider-nutanix/pull/147) ([JonKohler](https://github.com/JonKohler))
- WIP: Test cov \#141 [\#146](https://github.com/nutanix/terraform-provider-nutanix/pull/146) ([Crizstian](https://github.com/Crizstian))
- Change LICENSE to MPL2 per hashi [\#143](https://github.com/nutanix/terraform-provider-nutanix/pull/143) ([JonKohler](https://github.com/JonKohler))
- lint changes [\#142](https://github.com/nutanix/terraform-provider-nutanix/pull/142) ([JonKohler](https://github.com/JonKohler))
- squash lint issues and add unit tests via gotest [\#140](https://github.com/nutanix/terraform-provider-nutanix/pull/140) ([JonKohler](https://github.com/JonKohler))
- pulling in all develop work on travis/gnu [\#138](https://github.com/nutanix/terraform-provider-nutanix/pull/138) ([JonKohler](https://github.com/JonKohler))
- Master [\#137](https://github.com/nutanix/terraform-provider-nutanix/pull/137) ([JonKohler](https://github.com/JonKohler))
- Develop [\#136](https://github.com/nutanix/terraform-provider-nutanix/pull/136) ([JonKohler](https://github.com/JonKohler))
- Codecov \#112 Add Travis CI Code Coverage configuration [\#135](https://github.com/nutanix/terraform-provider-nutanix/pull/135) ([marinsalinas](https://github.com/marinsalinas))
- Image upload fail [\#134](https://github.com/nutanix/terraform-provider-nutanix/pull/134) ([thetonymaster](https://github.com/thetonymaster))
- remove more dead code [\#131](https://github.com/nutanix/terraform-provider-nutanix/pull/131) ([JonKohler](https://github.com/JonKohler))
- remove old code [\#130](https://github.com/nutanix/terraform-provider-nutanix/pull/130) ([JonKohler](https://github.com/JonKohler))
- one last line split fix [\#129](https://github.com/nutanix/terraform-provider-nutanix/pull/129) ([JonKohler](https://github.com/JonKohler))
- Pulling in test updates having challenges in GCP, but pass locally [\#128](https://github.com/nutanix/terraform-provider-nutanix/pull/128) ([JonKohler](https://github.com/JonKohler))
- squashing develop and master [\#126](https://github.com/nutanix/terraform-provider-nutanix/pull/126) ([JonKohler](https://github.com/JonKohler))
- fix quotes again [\#125](https://github.com/nutanix/terraform-provider-nutanix/pull/125) ([JonKohler](https://github.com/JonKohler))
- fix quotes [\#124](https://github.com/nutanix/terraform-provider-nutanix/pull/124) ([JonKohler](https://github.com/JonKohler))
- fix backticks [\#123](https://github.com/nutanix/terraform-provider-nutanix/pull/123) ([JonKohler](https://github.com/JonKohler))
- removing wait for ip and ip\_address from the vm ds and rs since not e… [\#121](https://github.com/nutanix/terraform-provider-nutanix/pull/121) ([Crizstian](https://github.com/Crizstian))
- updating example file and binaries to latest [\#120](https://github.com/nutanix/terraform-provider-nutanix/pull/120) ([Crizstian](https://github.com/Crizstian))
- Develop [\#118](https://github.com/nutanix/terraform-provider-nutanix/pull/118) ([Crizstian](https://github.com/Crizstian))
- Code refactor [\#117](https://github.com/nutanix/terraform-provider-nutanix/pull/117) ([marinsalinas](https://github.com/marinsalinas))
- Network sec rules refactor [\#116](https://github.com/nutanix/terraform-provider-nutanix/pull/116) ([marinsalinas](https://github.com/marinsalinas))
- Cloud init custom keys \#111 [\#113](https://github.com/nutanix/terraform-provider-nutanix/pull/113) ([marinsalinas](https://github.com/marinsalinas))
- Upload Image [\#110](https://github.com/nutanix/terraform-provider-nutanix/pull/110) ([thetonymaster](https://github.com/thetonymaster))
- Volume group rs 98 [\#109](https://github.com/nutanix/terraform-provider-nutanix/pull/109) ([marinsalinas](https://github.com/marinsalinas))
- Develop [\#108](https://github.com/nutanix/terraform-provider-nutanix/pull/108) ([Crizstian](https://github.com/Crizstian))
- refactorizing vm ds and rs [\#107](https://github.com/nutanix/terraform-provider-nutanix/pull/107) ([Crizstian](https://github.com/Crizstian))
- refactoring nsr ds and rs [\#106](https://github.com/nutanix/terraform-provider-nutanix/pull/106) ([Crizstian](https://github.com/Crizstian))
- Develop [\#105](https://github.com/nutanix/terraform-provider-nutanix/pull/105) ([Crizstian](https://github.com/Crizstian))
- Refactor [\#104](https://github.com/nutanix/terraform-provider-nutanix/pull/104) ([Crizstian](https://github.com/Crizstian))
- Develop [\#103](https://github.com/nutanix/terraform-provider-nutanix/pull/103) ([Crizstian](https://github.com/Crizstian))
- 97 subnets [\#102](https://github.com/nutanix/terraform-provider-nutanix/pull/102) ([Crizstian](https://github.com/Crizstian))
- 96 image ds [\#101](https://github.com/nutanix/terraform-provider-nutanix/pull/101) ([Crizstian](https://github.com/Crizstian))
- Documentation [\#95](https://github.com/nutanix/terraform-provider-nutanix/pull/95) ([marinsalinas](https://github.com/marinsalinas))
- Documentation [\#94](https://github.com/nutanix/terraform-provider-nutanix/pull/94) ([marinsalinas](https://github.com/marinsalinas))
- Refactor: client check errors [\#93](https://github.com/nutanix/terraform-provider-nutanix/pull/93) ([marinsalinas](https://github.com/marinsalinas))
- Network security rules \#80 [\#89](https://github.com/nutanix/terraform-provider-nutanix/pull/89) ([marinsalinas](https://github.com/marinsalinas))
- Develop [\#79](https://github.com/nutanix/terraform-provider-nutanix/pull/79) ([Crizstian](https://github.com/Crizstian))
- 68 cluster ds [\#78](https://github.com/nutanix/terraform-provider-nutanix/pull/78) ([Crizstian](https://github.com/Crizstian))
- Network Security Rule Data Source \#75 [\#77](https://github.com/nutanix/terraform-provider-nutanix/pull/77) ([marinsalinas](https://github.com/marinsalinas))
- Assign categories \#40 [\#76](https://github.com/nutanix/terraform-provider-nutanix/pull/76) ([marinsalinas](https://github.com/marinsalinas))
- Test passing \#71 [\#74](https://github.com/nutanix/terraform-provider-nutanix/pull/74) ([marinsalinas](https://github.com/marinsalinas))
- fix issue with vm.GuestCustomization not being initialized [\#73](https://github.com/nutanix/terraform-provider-nutanix/pull/73) ([htj](https://github.com/htj))
- Image test \#67 [\#70](https://github.com/nutanix/terraform-provider-nutanix/pull/70) ([marinsalinas](https://github.com/marinsalinas))
- fixes image rs [\#66](https://github.com/nutanix/terraform-provider-nutanix/pull/66) ([Crizstian](https://github.com/Crizstian))
- Categories resource \#36 [\#65](https://github.com/nutanix/terraform-provider-nutanix/pull/65) ([marinsalinas](https://github.com/marinsalinas))
- add completion callback for http request [\#64](https://github.com/nutanix/terraform-provider-nutanix/pull/64) ([marinsalinas](https://github.com/marinsalinas))
- Network Security Rule Resource \#30 [\#63](https://github.com/nutanix/terraform-provider-nutanix/pull/63) ([marinsalinas](https://github.com/marinsalinas))
- Add configuration for CI [\#58](https://github.com/nutanix/terraform-provider-nutanix/pull/58) ([marinsalinas](https://github.com/marinsalinas))
- Develop [\#55](https://github.com/nutanix/terraform-provider-nutanix/pull/55) ([JonKohler](https://github.com/JonKohler))
- Go fmt code [\#54](https://github.com/nutanix/terraform-provider-nutanix/pull/54) ([gliptak](https://github.com/gliptak))
- Develop [\#53](https://github.com/nutanix/terraform-provider-nutanix/pull/53) ([Crizstian](https://github.com/Crizstian))
- adding clusters data source to the provider, example.tf modified usin… [\#52](https://github.com/nutanix/terraform-provider-nutanix/pull/52) ([Crizstian](https://github.com/Crizstian))
- Go report issues \#47 [\#50](https://github.com/nutanix/terraform-provider-nutanix/pull/50) ([marinsalinas](https://github.com/marinsalinas))
- Develop [\#46](https://github.com/nutanix/terraform-provider-nutanix/pull/46) ([Crizstian](https://github.com/Crizstian))
- 25 cluster ds [\#45](https://github.com/nutanix/terraform-provider-nutanix/pull/45) ([Crizstian](https://github.com/Crizstian))
- Cleanup golint redundant nil check [\#44](https://github.com/nutanix/terraform-provider-nutanix/pull/44) ([gliptak](https://github.com/gliptak))
- Add GoReportCard badge [\#43](https://github.com/nutanix/terraform-provider-nutanix/pull/43) ([gliptak](https://github.com/gliptak))
- Add instructions for building project [\#42](https://github.com/nutanix/terraform-provider-nutanix/pull/42) ([marinsalinas](https://github.com/marinsalinas))
