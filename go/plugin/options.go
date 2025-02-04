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

func GetConsensusTarget(method *protogen.Method) *protonats.ConsensusTarget {
	if !proto.HasExtension(method.Desc.Options(), protonats.E_ConsensusTarget) {
		return nil
	}
	extension := proto.GetExtension(method.Desc.Options(), protonats.E_ConsensusTarget).(protonats.ConsensusTarget)
	return &extension
}

func IsConsensusLeader(method *protogen.Method) bool {
	if target := GetConsensusTarget(method); target == nil {
		return false
	} else {
		return *target == protonats.ConsensusTarget_LEADER
	}
}

func IsConsensusFollower(method *protogen.Method) bool {
	if target := GetConsensusTarget(method); target == nil {
		return false
	} else {
		return *target == protonats.ConsensusTarget_FOLLOWER
	}
}
