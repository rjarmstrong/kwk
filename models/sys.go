package models


type InfoResponse struct {
	Version         string `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
	Build           string `protobuf:"bytes,2,opt,name=build" json:"build,omitempty"`
	ClientSupported string `protobuf:"bytes,3,opt,name=clientSupported" json:"clientSupported,omitempty"`
	Dc              string `protobuf:"bytes,4,opt,name=dc" json:"dc,omitempty"`
}