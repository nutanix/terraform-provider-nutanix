package clustersv2

import (
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
)

// Define slices of enums
var (
	HypervisorTypes = []config.HypervisorType{
		config.HYPERVISORTYPE_AHV,
		config.HYPERVISORTYPE_ESX,
		config.HYPERVISORTYPE_HYPERV,
		config.HYPERVISORTYPE_XEN,
		config.HYPERVISORTYPE_NATIVEHOST,
	}

	SMTPTypes = []config.SmtpType{
		config.SMTPTYPE_PLAIN,
		config.SMTPTYPE_STARTTLS,
		config.SMTPTYPE_SSL,
	}

	ManagementServerTypes = []config.ManagementServerType{
		config.MANAGEMENTSERVERTYPE_VCENTER,
	}

	HTTPProxyTypes = []config.HttpProxyType{
		config.HTTPPROXYTYPE_HTTP,
		config.HTTPPROXYTYPE_HTTPS,
		config.HTTPPROXYTYPE_SOCKS,
	}

	HTTPProxyWhiteListTargetTypesSlice = []config.HttpProxyWhiteListTargetType{
		config.HTTPPROXYWHITELISTTARGETTYPE_IPV6_ADDRESS,
		config.HTTPPROXYWHITELISTTARGETTYPE_HOST_NAME,
		config.HTTPPROXYWHITELISTTARGETTYPE_DOMAIN_NAME_SUFFIX,
		config.HTTPPROXYWHITELISTTARGETTYPE_IPV4_NETWORK_MASK,
		config.HTTPPROXYWHITELISTTARGETTYPE_IPV4_ADDRESS,
	}

	ClusterFunctionRefs = []config.ClusterFunctionRef{
		config.CLUSTERFUNCTIONREF_AOS,
		config.CLUSTERFUNCTIONREF_ONE_NODE,
		config.CLUSTERFUNCTIONREF_TWO_NODE,
	}

	ClusterArchReferences = []config.ClusterArchReference{
		config.CLUSTERARCHREFERENCE_X86_64,
		config.CLUSTERARCHREFERENCE_PPC64LE,
	}

	DomainAwarenessLevels = []config.DomainAwarenessLevel{
		config.DOMAINAWARENESSLEVEL_RACK,
		config.DOMAINAWARENESSLEVEL_NODE,
		config.DOMAINAWARENESSLEVEL_BLOCK,
		config.DOMAINAWARENESSLEVEL_DISK,
	}

	ClusterFaultToleranceLevels = []config.ClusterFaultToleranceRef{
		config.CLUSTERFAULTTOLERANCEREF_CFT_1N_OR_1D,
		config.CLUSTERFAULTTOLERANCEREF_CFT_2N_OR_2D,
		config.CLUSTERFAULTTOLERANCEREF_CFT_1N_AND_1D,
		config.CLUSTERFAULTTOLERANCEREF_CFT_0N_AND_0D,
	}

	PIIScrubbingLevels = []config.PIIScrubbingLevel{
		config.PIISCRUBBINGLEVEL_ALL,
		config.PIISCRUBBINGLEVEL_DEFAULT,
	}

	KeyManagementServerTypes = []config.KeyManagementServerType{
		config.KEYMANAGEMENTSERVERTYPE_LOCAL,
		config.KEYMANAGEMENTSERVERTYPE_PRISM_CENTRAL,
		config.KEYMANAGEMENTSERVERTYPE_EXTERNAL,
	}

	OperationModes = []config.OperationMode{
		config.OPERATIONMODE_NORMAL,
		config.OPERATIONMODE_READ_ONLY,
		config.OPERATIONMODE_STAND_ALONE,
		config.OPERATIONMODE_SWITCH_TO_TWO_NODE,
		config.OPERATIONMODE_OVERRIDE,
	}

	EncryptionStatuses = []config.EncryptionStatus{
		config.ENCRYPTIONSTATUS_ENABLED,
		config.ENCRYPTIONSTATUS_DISABLED,
	}

	UpgradeStatuses = []config.UpgradeStatus{
		config.UPGRADESTATUS_PENDING,
		config.UPGRADESTATUS_DOWNLOADING,
		config.UPGRADESTATUS_QUEUED,
		config.UPGRADESTATUS_PREUPGRADE,
		config.UPGRADESTATUS_UPGRADING,
		config.UPGRADESTATUS_SUCCEEDED,
		config.UPGRADESTATUS_FAILED,
		config.UPGRADESTATUS_CANCELLED, //nolint:misspell
		config.UPGRADESTATUS_SCHEDULED,
	}
)

// Generate slices of string names for schema validation
var (
	HypervisorTypeStrings           = common.EnumToStrings(HypervisorTypes)
	SMTPTypeStrings                 = common.EnumToStrings(SMTPTypes)
	ManagementServerTypeStrings     = common.EnumToStrings(ManagementServerTypes)
	HTTPProxyTypeStrings            = common.EnumToStrings(HTTPProxyTypes)
	HTTPProxyWhiteListTargetStrings = common.EnumToStrings(HTTPProxyWhiteListTargetTypesSlice)
	ClusterFunctionStrings          = common.EnumToStrings(ClusterFunctionRefs)
	ClusterArchStrings              = common.EnumToStrings(ClusterArchReferences)
	DomainAwarenessLevelStrings     = common.EnumToStrings(DomainAwarenessLevels)
	ClusterFaultToleranceStrings    = common.EnumToStrings(ClusterFaultToleranceLevels)
	PIIScrubbingLevelStrings        = common.EnumToStrings(PIIScrubbingLevels)
	KeyManagementServerTypeStrings  = common.EnumToStrings(KeyManagementServerTypes)
	OperationModeStrings            = common.EnumToStrings(OperationModes)
	EncryptionStatusStrings         = common.EnumToStrings(EncryptionStatuses)
	UpgradeStatusStrings            = common.EnumToStrings(UpgradeStatuses)
)

// Generate maps of enum names to enum values for use in resource expansion
var (
	HypervisorTypeMap           = common.EnumToMap(HypervisorTypes)
	SMTPTypeMap                 = common.EnumToMap(SMTPTypes)
	ManagementServerTypeMap     = common.EnumToMap(ManagementServerTypes)
	HTTPProxyTypeMap            = common.EnumToMap(HTTPProxyTypes)
	HTTPProxyWhiteListTargetMap = common.EnumToMap(HTTPProxyWhiteListTargetTypesSlice)
	ClusterFunctionMap          = common.EnumToMap(ClusterFunctionRefs)
	ClusterArchMap              = common.EnumToMap(ClusterArchReferences)
	DomainAwarenessLevelMap     = common.EnumToMap(DomainAwarenessLevels)
	ClusterFaultToleranceMap    = common.EnumToMap(ClusterFaultToleranceLevels)
	PIIScrubbingLevelMap        = common.EnumToMap(PIIScrubbingLevels)
	KeyManagementServerTypeMap  = common.EnumToMap(KeyManagementServerTypes)
	OperationModeMap            = common.EnumToMap(OperationModes)
	EncryptionStatusMap         = common.EnumToMap(EncryptionStatuses)
	UpgradeStatusMap            = common.EnumToMap(UpgradeStatuses)
)
