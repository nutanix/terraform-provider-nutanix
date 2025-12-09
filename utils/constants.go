package utils

// RelEntityType constants - Relation Entity Types for the task entities affected
const (
	RelEntityTypeFloatingIPs            = "networking:config:floating-ips"
	RelEntityTypeVM                     = "vmm:ahv:config:vm"
	RelEntityTypeImages                 = "vmm:content:image"
	RelEntityTypeImagePlacementPolicy   = "vmm:images:config:placement-policy"
	RelEntityTypeTemplates              = "vmm:content:template"
	RelEntityTypeVolumeGroup            = "volumes:config:volume-group"
	RelEntityTypeVolumeGroupDisk        = "volumes:config:volume-group:disk"
	RelEntityTypeIscsiClient            = "volumes:config:iscsi-client"
	RelEntityTypeVPC                    = "networking:config:vpc"
	RelEntityTypeSubnet                 = "networking:config:subnet"
	RelEntityTypeFloatingIP             = "networking:config:floating-ip"
	RelEntityTypePBRS                   = "networking:config:routing-policy"
	RelEntityTypeSecurityPolicy         = "microseg:config:policy"
	RelEntityTypeServiceGroup           = "microseg:config:service-group"
	RelEntityTypeAddressGroup           = "microseg:config:address-group"
	RelEntityTypeVMDisk                 = "vmm:ahv:config:vm:disk"
	RelEntityTypeCDROM                  = "vmm:ahv:config:vm:cdrom"
	RelEntityTypeSerialPort             = "vmm:ahv:config:vm:serialport"
	RelEntityTypeVMNIC                  = "vmm:ahv:config:vm:nic"
	RelEntityTypeRecoveryPoint          = "dataprotection:config:recovery-point"
	RelEntityTypeVMRecoveryPoint        = "dataprotection:config:vm-recovery-point"
	RelEntityTypeStorageContainer       = "clustermgmt:config:storage-containers"
	RelEntityTypeRoute                  = "networking:config:route"
	RelEntityTypeObjects                = "objects:config:object-store"
	RelEntityTypeObjectStoreCertificate = "objects:config:object-store:certificate"
	RelEntityTypeOVA                    = "vmm:content:ova"
	RelEntityTypeStoragePolicy          = "datapolicies:config:storage-policy"
	RelEntityTypeKMS                    = "security:encryption:key-management-server"
	RelEntityTypeClusterProfile         = "clustermgmt:config:cluster-profile"
	RelEntityTypeDomainManager          = "prism:config:domain_manager"
)

// CompletionDetailsName constants - Completion details name for the task entities affected
const (
	CompletionDetailsNameRecoveryPoint    = "recoveryPointExtId"
	CompletionDetailsNameVMExtIDs         = "vmExtIds"
	CompletionDetailsNameVGExtIDs         = "volumeGroupExtIds"
	CompletionDetailsNameProtectionPolicy = "protectionPolicyExtId"
)
