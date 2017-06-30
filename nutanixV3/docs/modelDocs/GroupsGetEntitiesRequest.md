# GroupsGetEntitiesRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BucketBoundary** | **int64** | For grouping, the boundary to snap to when grouping. | [optional] [default to null]
**DownsamplingInterval** | **int64** | Downsampling interval to apply to query if override is desired.  | [optional] [default to null]
**EntityIds** | **[]string** | A set of entities that the request will be scoped to. | [optional] [default to null]
**EntityType** | **string** | The entity type that will be requested. | [default to null]
**FilterCriteria** | **string** | FIQL filter criteria that will be used to filter the returned data.  | [optional] [default to null]
**GroupAttributes** | [**[]GroupsRequestedAttribute**](groups_requested_attribute.md) |  | [optional] [default to null]
**GroupCount** | **int64** | The maximum number of groups to return in the result. | [optional] [default to null]
**GroupMemberAttributes** | [**[]GroupsRequestedAttribute**](groups_requested_attribute.md) |  | [optional] [default to null]
**GroupMemberCount** | **int64** | The maximum number of members to return per group. | [optional] [default to null]
**GroupMemberOffset** | **int64** | The offset into the total member set to return per group. | [optional] [default to null]
**GroupMemberSortAttribute** | **string** | The name of the attribute that will be used to sort group members.  | [optional] [default to null]
**GroupMemberSortDownsamplingFunction** | **string** | Downsampling function to take time series data and resolve to one value for sorting purposes.  | [optional] [default to null]
**GroupMemberSortOrder** | **string** | Sort order for entities and entity groups. | [optional] [default to null]
**GroupOffset** | **int64** | The offset into the total set of groups to return. | [optional] [default to null]
**GroupSortAttribute** | **string** | The name of the attribute that will be used to sort groups.  | [optional] [default to null]
**GroupSortDownsampleFunction** | **string** | Downsampling function to take time series data and resolve to one value for sorting purposes.  | [optional] [default to null]
**GroupSortOrder** | **string** | Sort order for entities and entity groups. | [optional] [default to null]
**GroupingAttribute** | **string** | Attribute that will be used to perform a group-by if needed.  | [optional] [default to null]
**GroupingAttributeType** | **string** | The type of an attribute being used for grouping - may be continuous or discrete.  | [optional] [default to null]
**IntervalEndMs** | **int64** | For a time-series query, the end of the interval since the epoch in ms. Default is latest value only.  | [optional] [default to 0]
**IntervalStartMs** | **int64** | For a time-series query, the start of the interval since the epoch in ms. Default is latest value only.  | [optional] [default to 0]
**NumberOfBuckets** | **int64** | For grouping, how many groups to return. | [optional] [default to null]
**NumberOfIntervalsForLatestData** | **int64** | When retrieving latest values, how far back to look as a multiple of the downsampling interval for the metric.  | [optional] [default to null]
**QueryName** | **string** | A custom name to use for tagging the query when debugging. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


