package clustersv2

import (
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
)

// ############################
// ###### Cluster Enums #######
// ############################
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

	PrivateKeyAlgorithms = []config.PrivateKeyAlgorithm{
		config.PRIVATEKEYALGORITHM_JKS,
		config.PRIVATEKEYALGORITHM_PKCS12,
		config.PRIVATEKEYALGORITHM_RSA_2048,
		config.PRIVATEKEYALGORITHM_RSA_4096,
		config.PRIVATEKEYALGORITHM_RSA_PUBLIC,
		config.PRIVATEKEYALGORITHM_KRB_KEYTAB,
		config.PRIVATEKEYALGORITHM_ECDSA_256,
		config.PRIVATEKEYALGORITHM_ECDSA_384,
		config.PRIVATEKEYALGORITHM_ECDSA_521,
	}
)

// ############################
// ###### Cluster Profiles ####
// ############################
var (
	AllowedOverrides = []config.ConfigType{
		config.CONFIGTYPE_NFS_SUBNET_WHITELIST_CONFIG,
		config.CONFIGTYPE_NTP_SERVER_CONFIG,
		config.CONFIGTYPE_SNMP_SERVER_CONFIG,
		config.CONFIGTYPE_SMTP_SERVER_CONFIG,
		config.CONFIGTYPE_PULSE_CONFIG,
		config.CONFIGTYPE_NAME_SERVER_CONFIG,
		config.CONFIGTYPE_RSYSLOG_SERVER_CONFIG,
	}
	SnmpAuthTypes = []config.SnmpAuthType{
		config.SNMPAUTHTYPE_MD5,
		config.SNMPAUTHTYPE_SHA,
	}
	SnmpPrivTypes = []config.SnmpPrivType{
		config.SNMPPRIVTYPE_DES,
		config.SNMPPRIVTYPE_AES,
	}
	SnmpProtocols = []config.SnmpProtocol{
		config.SNMPPROTOCOL_TCP,
		config.SNMPPROTOCOL_TCP6,
		config.SNMPPROTOCOL_UDP,
		config.SNMPPROTOCOL_UDP6,
	}
	SnmpTrapVersions = []config.SnmpTrapVersion{
		config.SNMPTRAPVERSION_V2,
		config.SNMPTRAPVERSION_V3,
	}
	RsyslogNetworkProtocols = []config.RsyslogNetworkProtocol{
		config.RSYSLOGNETWORKPROTOCOL_UDP,
		config.RSYSLOGNETWORKPROTOCOL_TCP,
		config.RSYSLOGNETWORKPROTOCOL_RELP,
	}
	RsyslogModuleNames = []config.RsyslogModuleName{
		config.RSYSLOGMODULENAME_AUDIT,
		config.RSYSLOGMODULENAME_CALM,
		config.RSYSLOGMODULENAME_MINERVA_CVM,
		config.RSYSLOGMODULENAME_STARGATE,
		config.RSYSLOGMODULENAME_FLOW_SERVICE_LOGS,
		config.RSYSLOGMODULENAME_SYSLOG_MODULE,
		config.RSYSLOGMODULENAME_CEREBRO,
		config.RSYSLOGMODULENAME_API_AUDIT,
		config.RSYSLOGMODULENAME_GENESIS,
		config.RSYSLOGMODULENAME_PRISM,
		config.RSYSLOGMODULENAME_ZOOKEEPER,
		config.RSYSLOGMODULENAME_FLOW,
		config.RSYSLOGMODULENAME_EPSILON,
		config.RSYSLOGMODULENAME_ACROPOLIS,
		config.RSYSLOGMODULENAME_UHARA,
		config.RSYSLOGMODULENAME_LCM,
		config.RSYSLOGMODULENAME_APLOS,
		config.RSYSLOGMODULENAME_NCM_AIOPS,
		config.RSYSLOGMODULENAME_CURATOR,
		config.RSYSLOGMODULENAME_CASSANDRA,
		config.RSYSLOGMODULENAME_LAZAN,
	}
	RsyslogLogSeverityLevels = []config.RsyslogModuleLogSeverityLevel{
		config.RSYSLOGMODULELOGSEVERITYLEVEL_EMERGENCY,
		config.RSYSLOGMODULELOGSEVERITYLEVEL_NOTICE,
		config.RSYSLOGMODULELOGSEVERITYLEVEL_ERROR,
		config.RSYSLOGMODULELOGSEVERITYLEVEL_ALERT,
		config.RSYSLOGMODULELOGSEVERITYLEVEL_INFO,
		config.RSYSLOGMODULELOGSEVERITYLEVEL_WARNING,
		config.RSYSLOGMODULELOGSEVERITYLEVEL_DEBUG,
		config.RSYSLOGMODULELOGSEVERITYLEVEL_CRITICAL,
	}
)

// ############################
// ###### Enum Helpers #########
// ############################
// Generate slices of string names for schema validation
var (
	// Cluster
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

	// Cluster Profiles
	AllowedOverridesStrings        = common.EnumToStrings(AllowedOverrides)
	SnmpAuthTypeStrings            = common.EnumToStrings(SnmpAuthTypes)
	SnmpPrivTypeStrings            = common.EnumToStrings(SnmpPrivTypes)
	SnmpProtocolStrings            = common.EnumToStrings(SnmpProtocols)
	SnmpTrapVersionStrings         = common.EnumToStrings(SnmpTrapVersions)
	RsyslogNetworkProtocolStrings  = common.EnumToStrings(RsyslogNetworkProtocols)
	RsyslogModuleNameStrings       = common.EnumToStrings(RsyslogModuleNames)
	RsyslogLogSeverityLevelStrings = common.EnumToStrings(RsyslogLogSeverityLevels)
	PrivateKeyAlgorithmStrings     = common.EnumToStrings(PrivateKeyAlgorithms)
)

// Generate maps of enum names to enum values for use in resource expansion
var (
	// Cluster
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

	UpgradeStatusMap = common.EnumToMap(UpgradeStatuses)

	// Cluster Profiles
	AllowedOverridesMap        = common.EnumToMap(AllowedOverrides)
	SnmpAuthTypeMap            = common.EnumToMap(SnmpAuthTypes)
	SnmpPrivTypeMap            = common.EnumToMap(SnmpPrivTypes)
	SnmpProtocolMap            = common.EnumToMap(SnmpProtocols)
	SnmpTrapVersionMap         = common.EnumToMap(SnmpTrapVersions)
	RsyslogNetworkProtocolMap  = common.EnumToMap(RsyslogNetworkProtocols)
	RsyslogModuleNameMap       = common.EnumToMap(RsyslogModuleNames)
	RsyslogLogSeverityLevelMap = common.EnumToMap(RsyslogLogSeverityLevels)
	PrivateKeyAlgorithmMap     = common.EnumToMap(PrivateKeyAlgorithms)
)
