package api

import (
	"gitlab2.psbc.com/ecn/ecn-xdsserver/common/constant"
	"gitlab2.psbc.com/ecn/ecn-xdsserver/xds"
	"sync"
)

type XdsGeneratorFactory struct {
	generatorMap map[string]XdsGenerator
}

var singletonXdsGeneratorFactory *XdsGeneratorFactory
var once sync.Once

func GetInstance() *XdsGeneratorFactory {
	once.Do(func() {
		_generatorMap := make(map[string]XdsGenerator)
		_generatorMap[constant.SERVICE_ENTRY_TYPE] = xds.GetInstance()
		// TODO Support other type generator
		singletonXdsGeneratorFactory = &XdsGeneratorFactory{generatorMap: _generatorMap}
	})
	return singletonXdsGeneratorFactory
}

func (factory *XdsGeneratorFactory) GetXdsGenerator(typeUrl string) XdsGenerator {
	xdsGenerator := factory.generatorMap[typeUrl]
	if xdsGenerator != nil {
		return xdsGenerator
	} else {
		return xds.GetSingletonEmptyXdsGenerator()
	}

}
