package v1alpha1

type ConditionType string

const (
	ConditionTypeReady      ConditionType = "Ready"
	ConditionTypeInProgress ConditionType = "InProgress"
	ConditionTypeFailed     ConditionType = "Failed"
	ConditionTypeUnknown    ConditionType = "Unknown"
)

type Condition struct {
	Type               ConditionType `json:"type,omitempty"`
	Status             bool          `json:"status,omitemtpy"`
	Reason             string        `json:"reason,omitempty"`
	Message            string        `json:"message,omitempty"`
	LastHeartbeatTime  string        `json:"lastHeartbeatTime,omitempty"`
	LastTransitionTime string        `json:"lastTransitionTime,omitempty"`
}
