package unversioned

type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

type ListMeta struct {
	SelfLink string `json:"selfLink,omitempty"`

	ResourceVersion string `json:"resourceVersion,omitempty"`

	TotalItems uint32 `json:"totalItems,omitempty"`
}
