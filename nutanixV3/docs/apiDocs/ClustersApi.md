#\ClustersApi

##ClustersListPost
//  Get a list of clusters
var api_instance := nutanix.ClustersApi
getEntitiesRequest := nutanix.ClusterListMetadata() // ClusterListMetadata | 
clusterslistpost_response, api_response, err := api_instance.ClustersListPost(getEntitiesRequest)

##ClustersUuidCertificatesCaCertsCaNameDelete
//  Delete the CA certificate
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
caName := "caName_example" // string | 
clustersuuidcertificatescacertscanamedelete_response, api_response, err := api_instance.ClustersUuidCertificatesCaCertsCaNameDelete(uuid, caName)

##ClustersUuidCertificatesCaCertsPost
//  Add a new CA certificate
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
spec := nutanix.CaCert() // CaCert | 
clustersuuidcertificatescacertspost_response, api_response, err := api_instance.ClustersUuidCertificatesCaCertsPost(uuid, spec)

##ClustersUuidCertificatesClientAuthDelete
//  Remove the CA chain for client authentication
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
clustersuuidcertificatesclientauthdelete_response, api_response, err := api_instance.ClustersUuidCertificatesClientAuthDelete(uuid)

##ClustersUuidCertificatesClientAuthPost
//  Import CA chain for client authentication
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
spec := nutanix.CaChainSpec() // CaChainSpec | 
clustersuuidcertificatesclientauthpost_response, api_response, err := api_instance.ClustersUuidCertificatesClientAuthPost(uuid, spec)

##ClustersUuidCertificatesClientAuthPut
//  Replace the CA chain for client authentication
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
spec := nutanix.CaChainSpec() // CaChainSpec | 
clustersuuidcertificatesclientauthput_response, api_response, err := api_instance.ClustersUuidCertificatesClientAuthPut(uuid, spec)

##ClustersUuidCertificatesCsrsGet
//  Download CSRs from cluster
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.CsrsSpec() // CsrsSpec | 
clustersuuidcertificatescsrsget_response, api_response, err := api_instance.ClustersUuidCertificatesCsrsGet(uuid, body)

##ClustersUuidCertificatesCsrsNodeIpGet
//  Download CSR from a discovered node
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
nodeIp := "nodeIp_example" // string | 
clustersuuidcertificatescsrsnodeipget_response, api_response, err := api_instance.ClustersUuidCertificatesCsrsNodeIpGet(uuid, nodeIp)

##ClustersUuidCertificatesPemkeyImportPost
//  Import a new key
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
spec := nutanix.PemkeySpec() // PemkeySpec | 
clustersuuidcertificatespemkeyimportpost_response, api_response, err := api_instance.ClustersUuidCertificatesPemkeyImportPost(uuid, spec)

##ClustersUuidCertificatesPemkeyPost
//  Generate a 2048 bits cipher length RSA key
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
clustersuuidcertificatespemkeypost_response, api_response, err := api_instance.ClustersUuidCertificatesPemkeyPost(uuid)

##ClustersUuidCertificatesSvmCertsKmsUuidPost
//  Add one or more certificates to a service VM
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
kmsUuid := "kmsUuid_example" // string | 
body := nutanix.CertificateSpecUploadInput() // CertificateSpecUploadInput | 
clustersuuidcertificatessvmcertskmsuuidpost_response, api_response, err := api_instance.ClustersUuidCertificatesSvmCertsKmsUuidPost(uuid, kmsUuid, body)

##ClustersUuidCertificatesSvmCertsNodeUuidKmsUuidDelete
//  Delete the certificate on a service VM
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
nodeUuid := "nodeUuid_example" // string | 
kmsUuid := "kmsUuid_example" // string | 
clustersuuidcertificatessvmcertsnodeuuidkmsuuiddelete_response, api_response, err := api_instance.ClustersUuidCertificatesSvmCertsNodeUuidKmsUuidDelete(uuid, nodeUuid, kmsUuid)

##ClustersUuidCertificatesSvmCertsNodeUuidKmsUuidPut
//  Replace the certificate on a service VM
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
nodeUuid := "nodeUuid_example" // string | 
kmsUuid := "kmsUuid_example" // string | 
cert := nutanix.Certificate() // Certificate | 
clustersuuidcertificatessvmcertsnodeuuidkmsuuidput_response, api_response, err := api_instance.ClustersUuidCertificatesSvmCertsNodeUuidKmsUuidPut(uuid, nodeUuid, kmsUuid, cert)

##ClustersUuidCloudCredentialsCloudTypeCloudCredentialsIdDelete
//  Delete a cloud credentials
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
cloudType := "cloudType_example" // string | 
cloudCredentialsId := 789 // int64 | 
clustersuuidcloudcredentialscloudtypecloudcredentialsiddelete_response, api_response, err := api_instance.ClustersUuidCloudCredentialsCloudTypeCloudCredentialsIdDelete(uuid, cloudType, cloudCredentialsId)

##ClustersUuidCloudCredentialsCloudTypeCloudCredentialsIdGet
//  Get a cloud credentials
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
cloudType := "cloudType_example" // string | 
cloudCredentialsId := 789 // int64 | 
clustersuuidcloudcredentialscloudtypecloudcredentialsidget_response, api_response, err := api_instance.ClustersUuidCloudCredentialsCloudTypeCloudCredentialsIdGet(uuid, cloudType, cloudCredentialsId)

##ClustersUuidCloudCredentialsCloudTypeCloudCredentialsIdPut
//  Update a cloud credentials
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
cloudType := "cloudType_example" // string | 
cloudCredentialsId := 789 // int64 | 
body := nutanix.CloudCredentialsIntentInput() // CloudCredentialsIntentInput | 
clustersuuidcloudcredentialscloudtypecloudcredentialsidput_response, api_response, err := api_instance.ClustersUuidCloudCredentialsCloudTypeCloudCredentialsIdPut(uuid, cloudType, cloudCredentialsId, body)

##ClustersUuidCloudCredentialsCloudTypeDelete
//  Delete all cloud credentials
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
cloudType := "cloudType_example" // string | 
clustersuuidcloudcredentialscloudtypedelete_response, api_response, err := api_instance.ClustersUuidCloudCredentialsCloudTypeDelete(uuid, cloudType)

##ClustersUuidCloudCredentialsCloudTypeListPost
//  Get a list of cloud credentials
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
cloudType := "cloudType_example" // string | 
getEntitiesRequest := nutanix.CloudCredentialsListMetadata() // CloudCredentialsListMetadata | 
clustersuuidcloudcredentialscloudtypelistpost_response, api_response, err := api_instance.ClustersUuidCloudCredentialsCloudTypeListPost(uuid, cloudType, getEntitiesRequest)

##ClustersUuidCloudCredentialsCloudTypePost
//  Add a cloud credential for accessing cloud sites
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
cloudType := "cloudType_example" // string | 
body := nutanix.CloudCredentialsIntentInput() // CloudCredentialsIntentInput | 
clustersuuidcloudcredentialscloudtypepost_response, api_response, err := api_instance.ClustersUuidCloudCredentialsCloudTypePost(uuid, cloudType, body)

##ClustersUuidGet
//  Get a cluster
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
clustersuuidget_response, api_response, err := api_instance.ClustersUuidGet(uuid)

##ClustersUuidPut
//  Update a cluster
var api_instance := nutanix.ClustersApi
uuid := "uuid_example" // string | The UUID of the entity
body := nutanix.ClusterIntentInput() // ClusterIntentInput | 
clustersuuidput_response, api_response, err := api_instance.ClustersUuidPut(uuid, body)

