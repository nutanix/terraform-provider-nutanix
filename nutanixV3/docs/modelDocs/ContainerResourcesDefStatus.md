# ContainerResourcesDefStatus

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ContainerOptions** | [**ContainerOptionsStatus**](container_options_status.md) | Various options for container. | [optional] [default to null]
**ContainerState** | **string** | Current state of the Container. | [optional] [default to null]
**ImageName** | **string** | Image name used for container. | [default to null]
**RegistryReference** | [**DockerRegistryReference**](docker_registry_reference.md) | Reference to container registry. | [optional] [default to null]
**RuntimeName** | **string** | Name of container given by runtime. | [optional] [default to null]
**SubnetReferenceList** | [**[]SubnetReference**](subnet_reference.md) | Networks associated with this container. | [optional] [default to null]
**VolumeList** | [**[]VolumeGroup**](volume_group.md) | Volumes associated with this container. | [optional] [default to null]
**VolumeReferenceList** | [**[]VolumeGroupReference**](volume_group_reference.md) | Referenced Volumes associated with this container. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


