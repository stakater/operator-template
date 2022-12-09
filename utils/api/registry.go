package api

type ObjectInformation struct {
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// +optional
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,6,opt,name=resourceVersion"`
}

type Registry = map[string]string
type IObjectTracking interface {
	GetRegistry() Registry
	Register(key string, object ObjectInformation)
	Unregister(key string)
}
