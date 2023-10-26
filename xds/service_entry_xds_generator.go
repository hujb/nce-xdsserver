package xds

import (
	"github.com/hujb/nce-xdsserver/common/resource"
	"google.golang.org/protobuf/types/known/anypb"
	"sync"
)

type ServiceEntryXdsGenerator struct {
}

var singletonServiceEntryXdsGenerator *ServiceEntryXdsGenerator

var once sync.Once

func GetInstance() *ServiceEntryXdsGenerator {
	once.Do(func() {
		singletonServiceEntryXdsGenerator = &ServiceEntryXdsGenerator{}
	})
	return singletonServiceEntryXdsGenerator
}

func (s *ServiceEntryXdsGenerator) Generate(rs *resource.ResourceSnapshot) []*anypb.Any {
	return nil
}
