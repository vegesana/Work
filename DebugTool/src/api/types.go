package api

type ObjectMeta struct {
	// Name of the object.
	Name string `json:"name"`
	// UUID of the object.
	UUID string `json:"uuid"`
	// SelfLink is a URL representing this object.
	SelfLink string `json:"selfLink, omitempty"`
	// Version of the resource. Can be used for change tracking.
	ResourceVersion   string            `json:"resourceVersion,omitempty"`
	Namespace         string            `json:"namespace,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	CreationTimestamp int64             `json:"creationTimestamp,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
}

type ObjectHealth string

const (
	// ObjectHealthGood - (Green) Service is active and no action is required.
	ObjectHealthGood ObjectHealth = "Good"
	// ObjectHealthWarning - (Orange) Service is unavailable.
	ObjectHealthWarning ObjectHealth = "Warning"
	// ObjectHealthFailed - (Red) Service is disrupted and action is required.
	ObjectHealthFailed ObjectHealth = "Failed"
	// ObjectHealthPending - (Yellow) Service is pending or in progress.
	ObjectHealthPending ObjectHealth = "Pending"
)
