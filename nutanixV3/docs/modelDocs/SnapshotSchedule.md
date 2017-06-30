# SnapshotSchedule

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DurationSecs** | **int64** | Duration of the event. If set, an event of duration duration_usecs will repeat as per the recurrence defined in interval_type  | [optional] [default to null]
**EndTime** | **int64** | End time of the snapshot schedule | [optional] [default to null]
**IntervalMultiple** | **int64** | Multiple of interval_type | [default to null]
**IntervalType** | **string** | Type of schedule interval | [default to null]
**IsSuspended** | **bool** | Whether the snapshot schedule is suspended | [optional] [default to null]
**StartTime** | **int64** | Start time of the snapshot schedule | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
