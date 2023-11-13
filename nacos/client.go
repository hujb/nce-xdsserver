package nacos

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/util"
	"github.com/nacos-group/nacos-sdk-go/vo"
	nceModel "github.com/nce/nce-xdsserver/model"
	"github.com/nce/nce-xdsserver/nacos/nacosDataProcess"
	"github.com/nce/nce-xdsserver/nacos/nacosResource"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"net/http"
	"time"
)

/*
*
通过客户端验证依赖包示例。
运行完后通过
curl -x "" 'http://127.0.0.1:8848/nacos/v1/ns/service/list?pageNo=1&pageSize=10&namespaceId=e525eafa-f7d7-4029-83d9-008937f9d468'
查看
*/
func main1() {
	sc := []constant.ServerConfig{
		{
			IpAddr: "192.168.1.83",
			Port:   8848,
		},
	}
	//or a more graceful way to create ServerConfig
	_ = []constant.ServerConfig{
		*constant.NewServerConfig("192.168.1.83", 8848),
	}

	cc := constant.ClientConfig{
		NamespaceId:         "test", //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogRollingConfig:    &lumberjack.Logger{MaxSize: 10},
		LogLevel:            "debug",
		AppendToStdout:      true,
	}
	//or a more graceful way to create ClientConfig
	_ = *constant.NewClientConfig(
		constant.WithNamespaceId("test"),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
		constant.WithLogRollingConfig(&lumberjack.Logger{MaxSize: 10}),
		constant.WithLogStdout(true),
	)

	// a more graceful way to create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(any(err))
	}

	//Register with default cluster and group
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   false,
		Metadata:    map[string]string{"idc": "shanghai"},
	})

	//Register with cluster name
	//GroupName=DEFAULT_GROUP
	ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.11",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-a",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})

	//Register different cluster
	//GroupName=DEFAULT_GROUP
	ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.12",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-b",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})

	//Register different group
	ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.13",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-b",
		GroupName:   "group-a",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.14",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-b",
		GroupName:   "group-b",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})

	//DeRegister with ip,port,serviceName
	//ClusterName=DEFAULT, GroupName=DEFAULT_GROUP
	//Note:ip=10.0.0.10,port=8848 should belong to the cluster of DEFAULT and the group of DEFAULT_GROUP.
	ExampleServiceClient_DeRegisterServiceInstance(client, vo.DeregisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		Ephemeral:   true, //it must be true
	})

	//DeRegister with ip,port,serviceName,cluster
	//GroupName=DEFAULT_GROUP
	//Note:ip=10.0.0.10,port=8848,cluster=cluster-a should belong to the group of DEFAULT_GROUP.
	ExampleServiceClient_DeRegisterServiceInstance(client, vo.DeregisterInstanceParam{
		Ip:          "10.0.0.11",
		Port:        8848,
		ServiceName: "demo.go",
		Cluster:     "cluster-a",
		Ephemeral:   true, //it must be true
	})

	//DeRegister with ip,port,serviceName,cluster,group
	ExampleServiceClient_DeRegisterServiceInstance(client, vo.DeregisterInstanceParam{
		Ip:          "10.0.0.14",
		Port:        8848,
		ServiceName: "demo.go",
		Cluster:     "cluster-b",
		GroupName:   "group-b",
		Ephemeral:   true, //it must be true
	})

	//Get service with serviceName
	//ClusterName=DEFAULT, GroupName=DEFAULT_GROUP
	ExampleServiceClient_GetService(client, vo.GetServiceParam{
		ServiceName: "demo.go",
	})
	//Get service with serviceName and cluster
	//GroupName=DEFAULT_GROUP
	ExampleServiceClient_GetService(client, vo.GetServiceParam{
		ServiceName: "demo.go",
		Clusters:    []string{"cluster-a", "cluster-b"},
	})
	//Get service with serviceName ,group
	//ClusterName=DEFAULT
	ExampleServiceClient_GetService(client, vo.GetServiceParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
	})

	//SelectAllInstance return all instances,include healthy=false,enable=false,weight<=0
	//ClusterName=DEFAULT, GroupName=DEFAULT_GROUP
	ExampleServiceClient_SelectAllInstances(client, vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
	})

	//SelectAllInstance
	//GroupName=DEFAULT_GROUP
	ExampleServiceClient_SelectAllInstances(client, vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
		Clusters:    []string{"cluster-a", "cluster-b"},
	})

	//SelectAllInstance
	//ClusterName=DEFAULT
	ExampleServiceClient_SelectAllInstances(client, vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
	})

	//SelectInstances only return the instances of healthy=${HealthyOnly},enable=true and weight>0
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	ExampleServiceClient_SelectInstances(client, vo.SelectInstancesParam{
		ServiceName: "demo.go",
	})

	//SelectOneHealthyInstance return one instance by WRR strategy for load balance
	//And the instance should be health=true,enable=true and weight>0
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	ExampleServiceClient_SelectOneHealthyInstance(client, vo.SelectOneHealthInstanceParam{
		ServiceName: "demo.go",
	})

	//Subscribe key=serviceName+groupName+cluster
	//Note:We call add multiple SubscribeCallback with the same key.
	param := &vo.SubscribeParam{
		ServiceName: "demo.go",
		Clusters:    []string{"cluster-b"},
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			fmt.Printf("callback111 return services:%s \n\n", util.ToJsonString(services))
		},
	}
	ExampleServiceClient_Subscribe(client, param)
	param2 := &vo.SubscribeParam{
		ServiceName: "demo.go",
		Clusters:    []string{"cluster-b"},
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			fmt.Printf("callback222 return services:%s \n\n", util.ToJsonString(services))
		},
	}
	ExampleServiceClient_Subscribe(client, param2)
	ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.112",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-b",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	//wait for client pull change from server
	time.Sleep(10 * time.Second)

	//Now we just unsubscribe callback1, and callback2 will still receive change event
	ExampleServiceClient_UnSubscribe(client, param)
	ExampleServiceClient_DeRegisterServiceInstance(client, vo.DeregisterInstanceParam{
		Ip:          "10.0.0.112",
		Ephemeral:   true,
		Port:        8848,
		ServiceName: "demo.go",
		Cluster:     "cluster-b",
	})
	//wait for client pull change from server
	time.Sleep(10 * time.Second)

	//GeAllService will get the list of service name
	//NameSpace default value is public.If the client set the namespaceId, NameSpace will use it.
	//GroupName default value is DEFAULT_GROUP
	ExampleServiceClient_GetAllService(client, vo.GetAllServiceInfoParam{
		PageNo:   1,
		PageSize: 10,
	})

	ExampleServiceClient_GetAllService(client, vo.GetAllServiceInfoParam{
		NameSpace: "0e83cc81-9d8c-4bb8-a28a-ff703187543f",
		PageNo:    1,
		PageSize:  10,
	})
}

func ExampleServiceClient_RegisterServiceInstance(client naming_client.INamingClient, param vo.RegisterInstanceParam) {
	success, _ := client.RegisterInstance(param)
	fmt.Printf("RegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func ExampleServiceClient_DeRegisterServiceInstance(client naming_client.INamingClient, param vo.DeregisterInstanceParam) {
	success, _ := client.DeregisterInstance(param)
	fmt.Printf("DeRegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func ExampleServiceClient_GetService(client naming_client.INamingClient, param vo.GetServiceParam) {
	service, _ := client.GetService(param)
	fmt.Printf("GetService,param:%+v, result:%+v \n\n", param, service)
}

func ExampleServiceClient_SelectAllInstances(client naming_client.INamingClient, param vo.SelectAllInstancesParam) {
	instances, _ := client.SelectAllInstances(param)
	fmt.Printf("SelectAllInstance,param:%+v, result:%+v \n\n", param, instances)
}

func ExampleServiceClient_SelectInstances(client naming_client.INamingClient, param vo.SelectInstancesParam) {
	instances, _ := client.SelectInstances(param)
	fmt.Printf("SelectInstances,param:%+v, result:%+v \n\n", param, instances)
}

func ExampleServiceClient_SelectOneHealthyInstance(client naming_client.INamingClient, param vo.SelectOneHealthInstanceParam) {
	instances, _ := client.SelectOneHealthyInstance(param)
	fmt.Printf("SelectInstances,param:%+v, result:%+v \n\n", param, instances)
}

func ExampleServiceClient_Subscribe(client naming_client.INamingClient, param *vo.SubscribeParam) {
	client.Subscribe(param)
}

func ExampleServiceClient_UnSubscribe(client naming_client.INamingClient, param *vo.SubscribeParam) {
	client.Unsubscribe(param)
}

func ExampleServiceClient_GetAllService(client naming_client.INamingClient, param vo.GetAllServiceInfoParam) {
	service, _ := client.GetAllServicesInfo(param)
	fmt.Printf("GetAllService,param:%+v, result:%+v \n\n", param, service)
}

// HttpGetNacosData 定义拉取nacos数据函数
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

func GetAllNamespaces(nacosUrl string) ([]string, error) {
	namespaceUrl := "http://" + nacosUrl + "/nacos/v1/console/namespaces"
	namespacestruct := &nacosDataProcess.NameSpaceMetadata{}

	nsmetadata, err := HttpGetNacosData(namespaceUrl)
	fmt.Printf("nsmetadata:%T\n", nsmetadata)
	if err != nil {
		log.Printf("获取namespaces列表失败，提示:%s", err)
		return nil, err
	}

	//将nsmetada转化为NameSpaceMetadata结构体类型
	err = json.Unmarshal(nsmetadata, namespacestruct)

	if err != nil {
		log.Printf("数据类型转化失败错误为：%v", err)
		return nil, err
	}
	var namespaces []string
	for _, nsdata := range namespacestruct.Data {
		namespaces = append(namespaces, nsdata.Namespace)
	}
	return namespaces, nil
}

func GetAllServicesByNamespace(nacosUrl string, param *nceModel.GetAllServiceInfoParam) ([]*nacosResource.NacosService, error) {
	//nacosUrl := serverConfig.IpAddr + ":" + string(serverConfig.Port)
	ServiceUrl := "http://" + nacosUrl + "/nacos/v1/ns/catalog/services?hasIpCount=true&withInstances=false&pageNo=1&pageSize=10&serviceNameParam=&groupNameParam=&namespaceId=" + param.NameSpace

	//获取service实例信息
	serviceMedata, err := HttpGetNacosData(ServiceUrl)
	if err != nil {
		return nil, errors.New("serviceMetadata 获取失败")
	}
	ServiceData := &nacosDataProcess.ServiceMetadata{}

	err = json.Unmarshal(serviceMedata, ServiceData)

	if err != nil {
		log.Printf("命名空间下的service数据获取失败,错误为：%v", err)
		return nil, errors.New("命名空间下的service数据获取失败,错误为 " + err.Error())
	}
	return ServiceData.ServiceList, nil
}

func GetAllServicesWithInstanceByNamespace(nacosUrl string, param *nceModel.GetAllServiceInfoParam) ([]*nacosResource.ServiceClusterInstanceDetail, error) {
	//nacosUrl := serverConfig.IpAddr + ":" + string(serverConfig.Port)
	ServiceUrl := "http://" + nacosUrl + "/nacos/v1/ns/catalog/services?hasIpCount=true&withInstances=true&pageNo=1&pageSize=10&serviceNameParam=&groupNameParam=&namespaceId=" + param.NameSpace

	//获取service实例信息
	serviceClusterInstanceMedata, err := HttpGetNacosData(ServiceUrl)
	if err != nil {
		log.Printf("serviceClusterInstanceMedata 获取失败,ServiceUrl=%s", ServiceUrl)
		panic(err)
		//return nil, errors.New("serviceClusterInstanceMedata 获取失败")
	}
	var ServiceClusterInstanceData []*nacosResource.ServiceClusterInstanceDetail

	err = json.Unmarshal(serviceClusterInstanceMedata, &ServiceClusterInstanceData)

	if err != nil {
		log.Printf("命名空间下的service数据获取失败,错误为：%v", err)
		panic(err)
		//return nil, err
	}

	return ServiceClusterInstanceData, nil
}

func GetNamespaceCheckSum(namespace string) string {
	return "1"
}
