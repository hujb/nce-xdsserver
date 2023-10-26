package api

import (
	"gitlab2.psbc.com/ecn/ecn-xdsserver/common/resource"
	"google.golang.org/protobuf/types/known/anypb"
)

type XdsGenerator interface {
	Generate(rs *resource.ResourceSnapshot) []*anypb.Any
}
