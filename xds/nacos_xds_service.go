package xds

import (
	"fmt"
	xds "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	_ "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/nce/nce-xdsserver/api"
	"github.com/nce/nce-xdsserver/common/constant"
	"github.com/nce/nce-xdsserver/common/resource"
	"github.com/nce/nce-xdsserver/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
	mcp_v1alpha1 "istio.io/api/mcp/v1alpha1"
	"istio.io/api/networking/v1alpha3"
	"strings"

	//"istio.io/api/security/v1beta1"
	security_v1beta1 "istio.io/api/security/v1beta1"
	"istio.io/client-go/pkg/apis/security/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"sync"
	"time"
)

type Connection struct {
	Stream xds.AggregatedDiscoveryService_StreamAggregatedResourcesServer
	// ConID is the connection identifier, used as a key in the connection table.
	// Currently based on the node name and a counter.
	ConID string

	NodeID string

	mu sync.RWMutex

	active bool

	// TODO PeerAddr is the address of the client envoy, from network layer
	PeerAddr string

	// TODO Metadata key-value pairs extending the Node identifier
	Metadata map[string]string

	NonceAcked map[string]string

	NonceSent map[string]string

	LastRequestTime int64

	LastRequestAcked bool
}

func (c Connection) String() string {
	return fmt.Sprintf("ConID: %s, NodeID: %s, Active: %t", c.ConID, c.NodeID, c.active)
}

type NacosXdsService struct {
	//pushc            chan struct{}
	clients          map[string]*Connection
	mutex            sync.RWMutex
	connectionNumber int
	resourceManager  *resource.ResourceManager
}

func NewNacosXdsService() *NacosXdsService {
	nacosMcpService := &NacosXdsService{
		clients: map[string]*Connection{},
	}
	return nacosMcpService
}

func (n *NacosXdsService) SetResourceManager(manager *resource.ResourceManager) {
	n.resourceManager = manager
}

func (n *NacosXdsService) StreamAggregatedResources(stream xds.AggregatedDiscoveryService_StreamAggregatedResourcesServer) error {
	log.Logger.Info("process ads stream.....")
	resource.GetResourceManagerInstance().InitResourceSnapshot()
	con := &Connection{
		Stream:           stream,
		LastRequestAcked: true,
		NonceAcked:       map[string]string{},
		NonceSent:        map[string]string{},
	}
	for {
		request, err := stream.Recv()
		//log.Printf("request info: %v", request.String())
		log.Logger.Debug("request info: " + request.String())
		if err != nil {
			if status.Code(err) == codes.Canceled || err == io.EOF {
				//log.Printf("ADS: %s terminated %v", con.ConID, err)
				log.Logger.Warn("ADS: " + con.ConID + "terminated " + err.Error())
				// remove this connection:
				delete(n.clients, con.ConID)
				return nil
			}
			//log.Printf("ADS: %s terminated with errors %v", con.ConID, err)
			log.Logger.Error("ADS: " + con.ConID + " terminated with errors " + err.Error())
			return err
		}
		err = n.Process(con, request, stream)
		if err != nil {
			return err
		}
	}
}

func (n *NacosXdsService) DeltaAggregatedResources(delta xds.AggregatedDiscoveryService_DeltaAggregatedResourcesServer) error {
	return nil
}

func (n *NacosXdsService) connectionID(node string) string {
	n.mutex.Lock()
	n.connectionNumber++
	c := n.connectionNumber
	n.mutex.Unlock()
	return node + "-" + strconv.Itoa(int(c))
}

func (n *NacosXdsService) Process(con *Connection, request *xds.DiscoveryRequest, stream xds.AggregatedDiscoveryService_StreamAggregatedResourcesServer) error {
	if !con.active {
		var id string
		if request.Node == nil || request.Node.Id == "" {
			log.Logger.Warn("Missing node id " + request.String())
			id = con.PeerAddr
		} else {
			id = request.Node.Id
		}

		con.mu.Lock()
		con.NodeID = id
		con.ConID = n.connectionID(con.NodeID)
		con.mu.Unlock()

		n.mutex.Lock()
		n.clients[con.ConID] = con
		n.mutex.Unlock()

		con.active = true

		//log.Println("activate new connection:", con)
		log.Logger.Info("activate new connection:" + con.String())
	}

	if !n.shouldPush(con, request) {
		return nil
	}

	if peerInfo, ok := peer.FromContext(stream.Context()); ok {
		log.Logger.Debug("peerInfo.Addr: " + peerInfo.Addr.String())
	}
	// 2023.10.30 22:57 start
	response := n.buildDiscoveryResponse(request.TypeUrl, n.resourceManager.GetResourceSnapshot())
	//log.Printf("DEBUG DiscoveryResponse: %v", response)
	log.Logger.Debug("DiscoveryResponse: " + response.String())
	con.NonceSent[response.TypeUrl] = response.Nonce
	err := stream.Send(response)
	// 2023.10.30 22:57 end

	//err := pushServiceEntries(request, con, stream) //2023.10.30 22:00废弃
	//err := pushPeerAuthentication(stream)  验证不同类型的资源
	if err != nil {
		// push failed - disconnect
		//log.Println("Closing connection ", con.ConID, err)
		log.Logger.Error("Closing connection " + con.ConID + ", err: " + err.Error())
		delete(n.clients, con.ConID)
		return err
	}

	return nil
}

func (n *NacosXdsService) shouldPush(con *Connection, request *xds.DiscoveryRequest) bool {
	rtype := request.TypeUrl

	if rtype == constant.MESH_CONFIG_TYPE {
		log.Logger.Warn("xds: type " + rtype + " should be ignored.")
		return false
	}

	if request.ErrorDetail != nil && request.ErrorDetail.Message != "" {
		var sb strings.Builder
		sb.WriteString("NACK: ")
		sb.WriteString(rtype)
		sb.WriteString(request.ErrorDetail.String())
		//log.Println("NACK: ", con.NodeID, rtype, request.ErrorDetail)
		log.Logger.Error(sb.String())
		return false
	}

	if request.ErrorDetail != nil && request.ErrorDetail.Code == 0 {
		con.mu.Lock()
		con.NonceAcked[rtype] = request.ResponseNonce
		con.mu.Unlock()
		log.Logger.Error("error: " + request.ErrorDetail.String())
		return false
	}

	if request.ResponseNonce != "" {
		// This shouldn't happen
		con.mu.Lock()
		lastNonce := con.NonceSent[rtype]
		con.mu.Unlock()

		if lastNonce == request.ResponseNonce {

			//if rtype == SERVICE_ENTRY_TYPE {
			var sb strings.Builder
			sb.WriteString("ACK of: ")
			sb.WriteString(strconv.Itoa(int(con.LastRequestTime)))
			sb.WriteString(", used time(microsecond): ")
			sb.WriteString(strconv.Itoa(int(time.Now().UnixNano()/1000 - con.LastRequestTime)))
			sb.WriteString("\n")
			log.Logger.Info(sb.String())
			//fmt.Println("ACK of:", con.LastRequestTime, " used time(microsecond): ", time.Now().UnixNano()/1000-con.LastRequestTime, "\n")
			con.LastRequestAcked = true

			con.mu.Lock()
			con.NonceAcked[rtype] = request.ResponseNonce
			con.mu.Unlock()
			//}

			return false
		} else {
			// will resent the resource, set the nonce - next response should be ok.
			//log.Println("Unmatching nonce ", request.ResponseNonce, lastNonce)
			log.Logger.Debug("Unmatching nonce " + request.ResponseNonce + ", lastNonce: " + lastNonce)
		}
	}
	return true
}

func pushServiceEntries(request *xds.DiscoveryRequest, con *Connection, stream xds.AggregatedDiscoveryService_StreamAggregatedResourcesServer) error {
	port := &v1alpha3.ServicePort{
		Number:   8080,
		Protocol: "HTTP",
		Name:     "http",
	}
	svcName := "demo.1"
	se := &v1alpha3.ServiceEntry{
		Hosts:      []string{svcName + ".nacos"},
		Resolution: v1alpha3.ServiceEntry_STATIC,
		Location:   v1alpha3.ServiceEntry_MESH_INTERNAL,
		Ports:      []*v1alpha3.ServicePort{port},
	}

	labels := make(map[string]string)
	labels["p"] = "hessian2"
	labels["ROUTE"] = "0"
	labels["APP"] = "ump"
	labels["st"] = "na62"
	labels["v"] = "2.0"
	labels["TIMEOUT"] = "3000"

	endpoint := &v1alpha3.WorkloadEntry{
		Labels: labels,
	}

	endpoint.Address = "0.0.0.1"
	endpoint.Ports = map[string]uint32{
		"http": uint32(8080),
	}

	se.Endpoints = append(se.Endpoints, endpoint)

	a, _ := anypb.New(se)
	//a, _ := types.MarshalAny(se)

	mcpResource := mcp_v1alpha1.Resource{
		Metadata: &mcp_v1alpha1.Metadata{
			Name:       "nacos/test",
			CreateTime: &timestamp.Timestamp{Seconds: time.Now().Unix()},
			Labels:     map[string]string{"hello": "test", "pa": "true"},
			Version:    "1111415485643",
		},
		Body: a,
	}

	apb, _ := anypb.New(&mcpResource)

	var _resources []*anypb.Any
	resources := append(_resources, apb)
	// Don't response a service entry query:  为什么会一直发networking.istio.io/v1alpha3/ServiceEntry请求，且node信息不为空
	// 是因为没有按照对应请求的typeUrl响应对应的资源？是的并且需要返回nil。没有收到ACK？ 是的。并且需要返回nil。（初步现象，深层次原理待研究）
	if request.TypeUrl != constant.SERVICE_ENTRY_TYPE {
		resources = []*anypb.Any{}
	}

	response := &xds.DiscoveryResponse{
		TypeUrl:     request.TypeUrl,
		VersionInfo: resource.GetResourceManagerInstance().GetResourceSnapshot().GetVersion(),
		Nonce:       fmt.Sprintf("%v", time.Now()),
		Resources:   resources,
	}
	//log.Printf("DEBUG DiscoveryResponse: %v", response)
	log.Logger.Debug("DiscoveryResponse: " + response.String())
	con.NonceSent[response.TypeUrl] = response.Nonce
	return stream.Send(response)
}

func pushPeerAuthentication(stream xds.AggregatedDiscoveryService_StreamAggregatedResourcesServer) error {
	pa := v1beta1.PeerAuthentication{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "security.istio.io/v1beta1",
			Kind:       "PeerAuthentication",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "istio-system",
		},
		Spec: security_v1beta1.PeerAuthentication{
			Mtls: &security_v1beta1.PeerAuthentication_MutualTLS{
				Mode: security_v1beta1.PeerAuthentication_MutualTLS_STRICT,
			},
		},
	}

	a, _ := anypb.New(&pa.Spec)

	mcpResource := mcp_v1alpha1.Resource{
		Metadata: &mcp_v1alpha1.Metadata{
			Name:       pa.ObjectMeta.Namespace + "/" + pa.ObjectMeta.Name,
			CreateTime: &timestamp.Timestamp{Seconds: time.Now().Unix()},
			Labels:     map[string]string{"hello": "test", "pa": "true"},
			Version:    "1111415485643",
		},
		Body: a,
	}

	apb, _ := anypb.New(&mcpResource)

	return stream.Send(&xds.DiscoveryResponse{
		TypeUrl:     "security.istio.io/v1beta1/PeerAuthentication",
		VersionInfo: "1",
		Nonce:       "",
		Resources:   []*anypb.Any{apb},
	})
}

// TODO 优化
func (n *NacosXdsService) HandleEvent(resourceSnapshot *resource.ResourceSnapshot) {
	if len(n.clients) == 0 {
		return
	}
	log.Logger.Info("xds: event changed trigger push.")
	serviceEntryResponse := n.buildDiscoveryResponse(constant.SERVICE_ENTRY_PROTO_PACKAGE, resourceSnapshot)
	log.Logger.Info("连接数：" + strconv.Itoa(len(n.clients)))
	for _, c := range n.clients {
		var sb strings.Builder
		sb.WriteString("sending resources count: ")
		sb.WriteString(strconv.Itoa(len(serviceEntryResponse.Resources)))
		sb.WriteString(", size: ")
		sb.WriteString(strconv.Itoa(int(n.sizeOfResources(serviceEntryResponse.Resources))))
		sb.WriteString(", request time: ")
		sb.WriteString(strconv.Itoa(int(c.LastRequestTime)))
		sb.WriteString(", connection id: ")
		sb.WriteString(c.ConID)
		log.Logger.Debug(sb.String())
		log.Logger.Debug("event changed DiscoveryResponse: " + serviceEntryResponse.String())
		//fmt.Println("sending resources count:", len(serviceEntryResponse.Resources), ", size:", n.sizeOfResources(serviceEntryResponse.Resources),
		//	", request time:", c.LastRequestTime, ", connection id:", c.ConID)
		c.NonceSent[serviceEntryResponse.TypeUrl] = serviceEntryResponse.Nonce
		c.Stream.Send(serviceEntryResponse)
	}

}

func (n *NacosXdsService) HandChangedEvent(resourceSnapshot *resource.ResourceSnapshot) {
	log.Logger.Info("xds: receive event changed trigger push.")
	if len(n.clients) == 0 {
		return
	}
	port := &v1alpha3.ServicePort{
		Number:   8081,
		Protocol: "HTTP",
		Name:     "http",
	}
	svcName := "demo.2"
	se := &v1alpha3.ServiceEntry{
		Hosts:      []string{svcName + ".nacos"},
		Resolution: v1alpha3.ServiceEntry_STATIC,
		Location:   v1alpha3.ServiceEntry_MESH_INTERNAL,
		Ports:      []*v1alpha3.ServicePort{port},
	}

	labels := make(map[string]string)
	labels["dc"] = "F"
	labels["DUS"] = "B001"
	labels["APP"] = "pcs"
	labels["TIMEOUT"] = "1000"

	endpoint := &v1alpha3.WorkloadEntry{
		Labels: labels,
	}

	endpoint.Address = "0.0.0.2"
	endpoint.Ports = map[string]uint32{
		"http": uint32(8081),
	}

	se.Endpoints = append(se.Endpoints, endpoint)

	a, _ := anypb.New(endpoint)
	//a, _ := types.MarshalAny(se)

	mcpResource := mcp_v1alpha1.Resource{
		Metadata: &mcp_v1alpha1.Metadata{
			Name:       "nacos/test2",
			CreateTime: &timestamp.Timestamp{Seconds: time.Now().Unix()},
			Labels:     map[string]string{"hello2": "test2", "pa2": "true2"},
			Version:    "1111415485644",
		},
		Body: a,
	}
	apb, _ := anypb.New(&mcpResource)
	log.Logger.Info("连接数：" + strconv.Itoa(len(n.clients)))
	for _, c := range n.clients {
		//if c.LastRequestAcked == false {
		//	log.Println("Last request not finished, ignore.")
		//	continue
		//}
		//c.LastRequestAcked = false

		c.LastRequestTime = time.Now().UnixNano() / 1000
		rs := []*anypb.Any{apb}
		var sb strings.Builder
		sb.WriteString("sending resources count: ")
		sb.WriteString(strconv.Itoa(len(rs)))
		sb.WriteString(", size: ")
		sb.WriteString(strconv.Itoa(int(n.sizeOfResources(rs))))
		sb.WriteString(", request time: ")
		sb.WriteString(strconv.Itoa(int(c.LastRequestTime)))
		sb.WriteString(", connection id: ")
		sb.WriteString(c.ConID)
		//log.Println("sending resources count:", len(rs), ", size:", n.sizeOfResources(rs),
		//	", request time:", c.LastRequestTime, ", connection id:", c.ConID)
		response := &xds.DiscoveryResponse{
			TypeUrl:     constant.SERVICE_ENTRY_TYPE,
			VersionInfo: resourceSnapshot.GetVersion(),
			Nonce:       fmt.Sprintf("%v", time.Now()),
			Resources:   rs,
		}
		//log.Printf("DEBUG event changed DiscoveryResponse: %v", response)
		log.Logger.Debug("event changed DiscoveryResponse: " + response.String())
		c.NonceSent[response.TypeUrl] = response.Nonce
		c.Stream.Send(response)
	}

}

func (n *NacosXdsService) sizeOfResources(rs []*anypb.Any) int64 {
	var length = 0
	for _, r := range rs {
		length = length + len(r.Value)
	}
	return int64(length)
}

func (n *NacosXdsService) buildDiscoveryResponse(typeUrl string, resourceSnapshot *resource.ResourceSnapshot) *xds.DiscoveryResponse {
	serviceEntryGenerator := api.GetInstance().GetXdsGenerator(typeUrl)
	rawResources := serviceEntryGenerator.Generate(resourceSnapshot)
	return &xds.DiscoveryResponse{
		TypeUrl:     typeUrl,
		VersionInfo: resourceSnapshot.GetVersion(),
		Nonce:       fmt.Sprintf("%v", time.Now()),
		Resources:   rawResources,
	}
}
