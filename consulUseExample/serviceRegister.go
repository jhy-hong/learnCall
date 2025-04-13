package consulUseExample

//import (
//	"fmt"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/health"
//	"google.golang.org/grpc/health/grpc_health_v1"
//	"log"
//	"net"
//)

/*
import (
	"fmt"
	"go-onepiece/common/consuleServiceDiscovery"
	"go-onepiece/onepieceService/hello/config"
	"go-onepiece/onepieceService/hello/internal"
	pb2 "go-onepiece/onepieceService/hello/pb"
	"go-onepiece/onepieceService/hello/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/gorm"
	"log"
	"net"
)

// 定义构造函数
func NewUserServer(db *gorm.DB) *service.HelloServer {
	return &service.HelloServer{
		DB: db, // 注入数据库连接
	}
}
func RegisterService() {
	log.Println("register service start")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.MyConfig.Grpc.Host, config.MyConfig.Grpc.Port))
	//lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// grpc的健康检查接口

	log.Println("register service to consul")
	// 将微服务注册到consul中  func RegisterServiceWithConsul(consulIp, serviceName, serviceIP string, consulPort, servicePort int)
	err = consuleServiceDiscovery.RegisterServiceWithConsul(config.MyConfig.Consul.Host, config.MyConfig.Grpc.Name, config.MyConfig.Grpc.Host, config.MyConfig.Consul.Port, config.MyConfig.Grpc.Port)
	if err != nil {
		log.Printf("failed to register service: %v\n", err)
		panic("failed to register service")
	}
	s := grpc.NewServer()

	// 初始化健康检查服务
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	// 设置服务状态为"SERVING"
	healthServer.SetServingStatus(config.MyConfig.Grpc.Name, grpc_health_v1.HealthCheckResponse_SERVING)
	pb2.RegisterHelloServiceServer(s, NewUserServer(internal.DB))

	// 启动服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


*/
