package xds

import (
	"github.com/hujb/nce-xdsserver/common/resource"
	"google.golang.org/protobuf/types/known/anypb"
	"sync"
)

type EmptyXdsGenerator struct {
}

var singletonEmptyXdsGenerator *EmptyXdsGenerator
var onceOnly sync.Once

func GetSingletonEmptyXdsGenerator() *EmptyXdsGenerator {
	onceOnly.Do(func() {
		singletonEmptyXdsGenerator = &EmptyXdsGenerator{}
	})
	return singletonEmptyXdsGenerator
}

func (e *EmptyXdsGenerator) Generate(rs *resource.ResourceSnapshot) []*anypb.Any {
	return nil
}
