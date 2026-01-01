package plugin

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"xiam.li/protonats/go/protonats"
)

func IsUsingBroadcasting(method *protogen.Method) bool {
	if !proto.HasExtension(method.Desc.Options(), protonats.E_Broadcast) {
		return false
	}
	extension := proto.GetExtension(method.Desc.Options(), protonats.E_Broadcast)
	return extension.(bool)
}
