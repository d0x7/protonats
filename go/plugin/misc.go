package plugin

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
)

func ProtocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

func SubjectName(service *protogen.Service, method *protogen.Method) string {
	return "service." + service.GoName + "." + method.GoName
}
