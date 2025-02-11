type  struct {
	Mode                      *string           `json:"mode,omitempty"`
	From                      *string           `json:"from,omitempty"`
	To                        *string           `json:"to,omitempty"`
	TimeUnit                  *string           `json:"timeUnit,omitempty"`
	TimeUnitNumber            *string           `json:"timeUnitNumber,omitempty"`
	DatabaseIds               []*string         `json:"databaseIds,omitempty"`
	Snapshots                 *ListSnapshots    `json:"snapshots,omitempty"`
	ContinuousRegion          *ContinuousRegion `json:"continuousRegion,omitempty"`
	DatabasesContinuousRegion interface{}       `json:"databasesContinuousRegion,omitempty"`
}