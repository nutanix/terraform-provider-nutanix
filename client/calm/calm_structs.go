package calm

import "time"

// Reference ...
type Reference struct {
	Kind *string `json:"kind" mapstructure:"kind"`
	Name *string `json:"name,omitempty" mapstructure:"name,omitempty"`
	UUID *string `json:"uuid" mapstructure:"uuid"`
}

// Metadata Metadata The kind metadata
type Metadata struct {
	LastUpdateTime   *time.Time        `json:"last_update_time,omitempty" mapstructure:"last_update_time,omitempty"`   //
	Kind             *string           `json:"kind" mapstructure:"kind"`                                               //
	UUID             *string           `json:"uuid,omitempty" mapstructure:"uuid,omitempty"`                           //
	ProjectReference *Reference        `json:"project_reference,omitempty" mapstructure:"project_reference,omitempty"` // project reference
	CreationTime     *time.Time        `json:"creation_time,omitempty" mapstructure:"creation_time,omitempty"`
	SpecVersion      *int64            `json:"spec_version,omitempty" mapstructure:"spec_version,omitempty"`
	SpecHash         *string           `json:"spec_hash,omitempty" mapstructure:"spec_hash,omitempty"`
	OwnerReference   *Reference        `json:"owner_reference,omitempty" mapstructure:"owner_reference,omitempty"`
	Categories       map[string]string `json:"categories,omitempty" mapstructure:"categories,omitempty"`
	Name             *string           `json:"name,omitempty" mapstructure:"name,omitempty"`

	// Applied on Prism Central only. Indicate whether force to translate the spec of the fanout request to fit the target cluster API schema.
	ShouldForceTranslate *bool `json:"should_force_translate,omitempty" mapstructure:"should_force_translate,omitempty"`

	//TODO: add if necessary
	//CategoriesMapping    map[string][]string `json:"categories_mapping,omitempty" mapstructure:"categories_mapping,omitempty"`
	//EntityVersion        *string             `json:"entity_version,omitempty" mapstructure:"entity_version,omitempty"`
	//UseCategoriesMapping *bool               `json:"use_categories_mapping,omitempty" mapstructure:"use_categories_mapping,omitempty"`
}

// MessageResource ...
type MessageResource struct {

	// Custom key-value details relevant to the status.
	Details map[string]string `json:"details,omitempty" mapstructure:"details,omitempty"`

	// If state is ERROR, a message describing the error.
	Message *string `json:"message" mapstructure:"message"`

	// If state is ERROR, a machine-readable snake-cased *string.
	Reason *string `json:"reason" mapstructure:"reason"`
}

type ProjectQuotaMetadata struct {
	Kind             *string    `json:"kind,omitempty"`
	ProjectReference *Reference `json:"project_reference,omitempty"`
	UUID             *string    `json:"uuid,omitempty"`
}

type ProjectQuotaData struct {
	Disk   *int `json:"disk,omitempty"`
	VCPU   *int `json:"vcpu,omitempty"`
	Memory *int `json:"memory,omitempty"`
}

type ProjectQuotaEntities struct {
	Project *string `json:"project,omitempty"`
}

type ProjectQuotaResources struct {
	Data       *ProjectQuotaData     `json:"data,omitempty"`
	Entities   *ProjectQuotaEntities `json:"entities,omitempty"`
	Metadata   *Metadata             `json:"metadata,omitempty"`
	UUID       *string               `json:"uuid,omitempty"`
	State      *string               `json:"state,omitempty"`
	EntityType *string               `json:"entity_type,omitempty"`
}

type ProjectQuotaSpec struct {
	Resources *ProjectQuotaResources `json:"resources,omitempty"`
}

// Project CALM Quota
type ProjectQuotaIntentInput struct {
	Metadata *ProjectQuotaMetadata `json:"metadata,omitempty"`
	Spec     *ProjectQuotaSpec     `json:"spec,omitempty"`
}

type ProjectQuotaStatus struct {
	State       *string                `json:"state,omitempty"`
	MessageList []*MessageResource     `json:"message_list,omitempty"`
	Resources   *ProjectQuotaResources `json:"resources,omitempty"`
}

type ProjectQuotaIntentResponse struct {
	Status   *ProjectQuotaStatus `json:"status,omitempty"`
	Metadata *Metadata           `json:"metadata,omitempty"`
	Spec     *ProjectQuotaSpec   `json:"spec,omitempty"`
}

type EnableProjectQuotaInput struct {
	Spec *ProjectQuotaSpec `json:"spec,omitempty"`
}
