package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/nce/nce-xdsserver/common/event/process"
	r "github.com/nce/nce-xdsserver/common/resource"
	rs "github.com/nce/nce-xdsserver/common/resource"
	nceModel "github.com/nce/nce-xdsserver/model"
	"github.com/nce/nce-xdsserver/nacos"
	"github.com/nce/nce-xdsserver/xds"
	"go.uber.org/zap"
	"istio.io/api/meta/v1alpha1"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	istio "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {

	//zapLogTest()
	//testGoLog()
	//testVirtualService()
	testStringsBuilder()
}

func testStringsBuilder() {
	var builder strings.Builder
	builder.WriteString("Hello")
	builder.WriteString(" ")
	builder.WriteString("World!")
	result := builder.String()
	fmt.Println(result) // 输出: Hello World!

}

func testVirtualService() {
	// 获取kubeconfig文件路径
	kubeconfig := flag.String("kubeconfig", "/Users/hujiebin/source/istio/kubeconfig1", "absolute path to the kubeconfig file")
	flag.Parse()
	// 创建Kubernetes和Istio的客户端集
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	//kubernetesClient, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	panic(err.Error())
	//}
	istioClient, err := istio.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
		//panic(err.Error())
	}

	// 创建VirtualService对象
	httpRoutes := make([]*networkingv1alpha3.HTTPRoute, 0)
	httpRoute := &networkingv1alpha3.HTTPRoute{
		Match: []*networkingv1alpha3.HTTPMatchRequest{
			{
				Uri: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Prefix{
						Prefix: "/",
					},
				},
			},
		},
		Route: []*networkingv1alpha3.HTTPRouteDestination{
			{
				Destination: &networkingv1alpha3.Destination{
					Host: "reviews",
					Port: &networkingv1alpha3.PortSelector{
						Number: 9091,
					},
				},
			},
		},
	}
	httpRoutes = append(httpRoutes, httpRoute)
	virtualService := &v1alpha3.VirtualService{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{Name: "my-virtual-service"},
		Spec: networkingv1alpha3.VirtualService{Hosts: []string{"my-virtual-service"},
			Http: httpRoutes, ExportTo: []string{"*"}},
		Status: v1alpha1.IstioStatus{},
	}

	// 应用VirtualService到Kubernetes集群
	createdVirtualService, err := istioClient.NetworkingV1alpha3().VirtualServices("default").Create(context.TODO(), virtualService, v1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created VirtualService: %s\n", createdVirtualService.ObjectMeta.Name)
}

func testGoLog() {
	logFile, err := os.Create("./" + time.Now().Format("20060102") + ".txt")
	if err != nil {
		fmt.Println(err)
	}
	//创建一个Logger
	//参数1：日志写入目的地
	//参数2：每条日志的前缀
	//参数3：日志属性
	logger := log.New(logFile, "test_", log.Ldate|log.Ltime)
	//Flags返回Logger的输出选项
	fmt.Println(logger.Flags())
	//SetFlags设置输出选项
	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	fmt.Println(logger.Prefix())
	//设置输出前缀
	logger.SetPrefix("test_")
	logger.Output(2, "打印一条日志信息")
}

func testVariable() {
	a := true
	b := a
	a = false
	fmt.Println(a)
	fmt.Println(b)
}

func testSnapshotVersion() {
	nacosXdsService := xds.NewNacosXdsService()
	eventProcess, _ := process.NewEventProcessor(nacosXdsService)
	watcher, _ := nacos.NewNacosServiceInfoResourceWatcher(eventProcess)
	r.GetResourceManagerInstance().SetResourceWatcher(watcher)
	rs1 := &rs.ResourceSnapshot{}
	rs1.InitResourceSnapshot(r.GetResourceManagerInstance())
	fmt.Println(rs1.GetVersion())

	rs2 := &rs.ResourceSnapshot{}
	rs2.InitResourceSnapshot(r.GetResourceManagerInstance())
	fmt.Println(rs2.GetVersion())
}

func testPoint() {
	//m := make(map[string]Shape)
	//m["a"] = &Rectangle{width: 1.0, height: 2.0}
	//r := &Rectangle{width: 1.0, height: 3.0}
	r1 := new(Rectangle)
	r1.height = 3.0
	r1.width = 1.0
	//log.Println("正方形修改前：", r1)
	//fmt.Println("面积=", r.Area())
	r1.Area()
	//log.Println("正方形修改后：", r1)
}

func testSyncMap() {
	m := sync.Map{}
	m.Store("s1", &nceModel.IstioService{})
	m.Store("s2", &nceModel.IstioService{})
	m.Range(func(key, value interface{}) bool {
		//fmt.Printf("key: %s, value: %v", key, value)
		return true
	})
}

func testChannel() {
	numjobs := 10
	jobs := make(chan int, numjobs)
	results := make(chan int, numjobs)
	go worker(1, jobs, results)
	go worker(2, jobs, results)
	for j := 1; j <= numjobs; j++ {
		jobs <- j
	}
	close(jobs)
	for i := 1; i <= numjobs; i++ {
		fmt.Println("Result:", <-results)
	}
	close(results)
}

func testGetNacosData() {
	nacosUrl := "127.0.0.1:8848"
	//ServiceEntry, _ := nacosDataProcess.ConstServiceEntry(nacosUrl)
	//fmt.Println("ServiceEntry=", ServiceEntry)
	// 测试获取所有nacos所有ns
	namespaces, err := nacos.GetAllNamespaces(nacosUrl)
	if err != nil {
		log.Print(err)
		return
	}
	for _, namespace := range namespaces {
		log.Printf("namespace=%s", namespace)
	}

	// 测试获取nacos指定ns下的所有instance
	param := &nceModel.QueryAllServiceInfoParam{NamespaceId: "public"}
	instances, err := nacos.GetAllServicesWithInstanceByNamespace(nacosUrl, param)
	if err != nil {
		log.Print(err)
		return
	}
	for _, instance := range instances {
		log.Printf("instance=%v", instance)
	}

	//测试获取nacos指定ns下的所有service
	services, err := nacos.GetAllServicesByNamespace(nacosUrl, param)
	if err != nil {
		log.Print(err)
		return
	}
	for _, service := range services {
		log.Printf("service=%v", service)
	}
}

type Shape interface {
	Area() float32
}

func getRectangle() (rec *Rectangle, err error) {
	return nil, errors.New("aaa")
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}

}

func zapLogTest() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	url := "http://localhost/"
	logger.Info("production failed to fetch URL", zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second))
}
