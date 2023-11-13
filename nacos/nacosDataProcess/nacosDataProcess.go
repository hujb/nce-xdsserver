package nacosDataProcess

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nce/nce-xdsserver/model"
	"github.com/nce/nce-xdsserver/nacos/nacosResource"
	"io"
	"log"
	"net/http"
)

// 定义拉取nacos数据函数
func HttpGetNacosData(nacosUrl string) (nacosData []byte, err error) {
	resp, err := http.Get(nacosUrl)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()
	nacosData, err = io.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.New("拉取nacos数据失败")
	}
	return nacosData, nil
}

// 定义请求nacos数据url类型
type RequestNacosUrl struct {
	GetNameSpaceMetadata *NameSpaceMetadata
	GetServiceMetadata   *ServiceMetadata
	GetInstanceMetadata  *InstanceMetadata
}

type NameSpaceMetadata struct {
	Code   int32        `json:"code"`
	Mssage string       `json:"message"`
	Data   []DataMetada `json:"data"`
}

type DataMetada struct {
	Namespace         string `json:"namespace"`
	NamespaceShowName string `json:"namespaceShowName"`
	NamespaceDesc     string `json:"namespaceDesc"`
	Quota             int32  `json:"quota"`
	Configcount       int32  `json:"configcount"`
	Type              int32  `json:"type"`
}

type ServiceMetadata struct {
	ServiceList []*nacosResource.NacosService `json:"serviceList"`
	Count       int32                         `json:"count"`
}

type InstanceMetadata struct {
	list  []map[string]string
	count int32
}

func ConstServiceEntry(nacosUrl string) (se model.ServiceEntry, err error) {
	//port := istioDataType.Port{
	//	Number:   8080,
	//	Protocol: "HTTP",
	//	Name:     "http",
	//}

	se = model.ServiceEntry{
		Ports:      make([]*model.Port, 0),
		Location:   model.ServiceEntry_MESH_INTERNAL,
		Resolution: model.ServiceEntry_STATIC,
	}
	fmt.Println("ServiceEntry:", se)

	//获取nacos元数据逐级构造serviceEntry
	namespaceUrl := "http://" + nacosUrl + "/nacos/v1/console/namespaces"
	namespacestruct := &NameSpaceMetadata{}

	nsmetadata, err := HttpGetNacosData(namespaceUrl)
	fmt.Printf("nsmetadata:%T\n", nsmetadata)
	if err != nil {
		log.Println("获取namespaces列表失败，提示%s", err)
		return se, err
	}

	//nsmetadataToString := string(nsmetadata)
	//fmt.Printf("nsmetadataToString=%T", nsmetadataToString)

	//将nsmetada转化为NameSpaceMetadata结构体类型
	err = json.Unmarshal(nsmetadata, namespacestruct)

	if err != nil {
		log.Printf("数据类型转化失败错误为：%v", err)
	}
	fmt.Println("namespacestruct=", namespacestruct)

	//接收循环转换后得到的serviceEntry数据存储到切片中
	//TotalServiceEntry := &[]istioDataType.ServiceEntry{}

	for _, nsdata := range namespacestruct.Data {

		//构造service Url获取nacos service信息
		ServiceUrl := "http://" + nacosUrl + "/nacos/v1/ns/catalog/services?hasIpCount=true&withInstances=false&pageNo=1&pageSize=10&serviceNameParam=&groupNameParam=&namespaceId=" + nsdata.NamespaceShowName

		//获取service实例信息
		serviceMedata, err := HttpGetNacosData(ServiceUrl)
		if err != nil {
			return se, errors.New("serviceMetadata 获取失败")
		}

		//将service转化为string
		//serviceMedataTostring := string(serviceMedata)
		//fmt.Printf("nsdata=%v\nserviceMedata=%v\n", nsdata.NamespaceShowName, serviceMedataTostring)

		//遍历nacos namespace后将得到的service元数据转化为事先定义好的ServiceMetadata数据类型
		ServiceData := &ServiceMetadata{}

		err = json.Unmarshal(serviceMedata, ServiceData)

		fmt.Println("ServiceData=", ServiceData)
		fmt.Println("Service=", ServiceData.ServiceList)

		//判断如果获取service失败则不在执行下面的逻辑
		if err != nil {
			log.Println("命名空间下的service数据获取失败,错误为：", err)
			continue
		}

		//转化为serviceMetadata结构体后，如果没有service信息也不再执行下面的逻辑
		if ServiceData.Count <= 0 {
			log.Println("该命名空间下未查询到service信息，可能没有instance实例注册，请确认！命名空间为：", nsdata.NamespaceShowName)
			continue
		}

		//转换nacos service为serviceEntry类型

		serviceSlice := ServiceData.ServiceList

		for _, svcValue := range serviceSlice {
			//构造serviceEntry
			//ServiceToServiceEntry := istioDataType.ServiceEntry{
			//	Ports:      make([]*istioDataType.Port, 0),
			//	Location:   istioDataType.ServiceEntry_MESH_INTERNAL,
			//	Resolution: istioDataType.ServiceEntry_STATIC,
			//}
			se.Hosts = append(se.Hosts, svcValue.Name+".ecn-nacos")
			fmt.Println("serviceEntry=", se)

			//请求service实例信息，获取cluster name
			ServiceDetailUrl := "http://" + nacosUrl + "/nacos/v1/ns/catalog/service?&groupName=&pageSize=10000&pageNo=1&namespaceId=" + nsdata.NamespaceShowName + "&serviceName=" + svcValue.Name

			ServiceDetail, err := HttpGetNacosData(ServiceDetailUrl)

			if err != nil {
				log.Println("获取service详情信息异常，错误提示：", err)
				continue
			}

			//将service详情信息转化为自定义数据结构
			SvcDetailStore := &nacosResource.ServiceClusterDetail{}
			SvcDetailStoreToString := string(ServiceDetail)
			fmt.Println("SvcDetailStoreToString=", SvcDetailStoreToString)

			err = json.Unmarshal(ServiceDetail, SvcDetailStore)

			if err != nil {
				log.Println("service 详情信息数据转换失败，错误提示：", err)
				continue
			}

			fmt.Println("SvcDetailStore=", SvcDetailStore)

			//获取cluster name列表
			for _, Cluster := range SvcDetailStore.Clusters {
				fmt.Println("Cluster=", Cluster)

				//构造workloadEntry
				InstanceUrl := "http://" + nacosUrl + "/nacos/v1/ns/catalog/instances?&pageSize=10000&pageNo=1&namespaceId=" + nsdata.NamespaceShowName + "&groupName=&serviceName=" + svcValue.Name + "&clusterName=" + Cluster.Name
				InstanceMetaData, err := HttpGetNacosData(InstanceUrl)
				InstanceMetaDataToString := string(InstanceMetaData)
				fmt.Println("InstanceMetaDataToString", InstanceMetaDataToString)

				if err != nil {
					log.Println("获取instance实例信息异常，报错提示：", err)
					continue
				}
				//将instance实例信息转换为结构体存储
				InstanceMedataStruct := &nacosResource.InstanceMetadata{}

				err = json.Unmarshal(InstanceMetaData, InstanceMedataStruct)
				if err != nil {
					log.Println("Instance数据转换异常，错误提示=", err)
				}

				fmt.Println("InstanceMetaDataStruct=", InstanceMedataStruct)

				//获取instance信息构造workloadEntry
				for _, Instances := range InstanceMedataStruct.List {
					//	构造workloadEntry结构体
					CunstWorkLoad := &model.ServiceEntry_Endpoint{
						Address: Instances.Ip,
						Ports:   Instances.Port,
						Labels:  Instances.Metadata,
						Weight:  uint32(Instances.Weight),
					}
					fmt.Println("CunstWorkLoad", CunstWorkLoad)

					//将生成的ServiceEntry_Endpoint添加到ServiceEntry中
					se.Endpoints = append(se.Endpoints, CunstWorkLoad)
				}
			}
		}
	}
	return se, err
}

//keyvalues := namespacestruct.data

//将map类型的namespace数据转换为
//for key, _ := range keyvalues {
//	switch key {
//	case name:
//		namespacestruct.data[name]
//
//	}
//}
