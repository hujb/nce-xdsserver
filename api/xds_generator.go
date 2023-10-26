package api

import (
	"github.com/hujb/nce-xdsserver/common/resource"
	"google.golang.org/protobuf/types/known/anypb"
)

type XdsGenerator interface {
	Generate(rs *resource.ResourceSnapshot) []*anypb.Any
}
