package consulUseExample

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"log"
	"math/rand"
)

func RegisterServiceWithConsul(consulIp, serviceName, serviceIP string, consulPort, servicePort int) error {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", consulIp, consulPort) // Consul地址
	client, err := api.NewClient(config)
	if err != nil {
		log.Printf("Consul connect fail,err:%v", err)
		return err
	}

	// 获取本机IP
	//addr, err := getLocalIP()
	//if err != nil {
	//	panic(err)
	//}
	id := uuid.New().String()

	// 服务注册信息
	registration := &api.AgentServiceRegistration{
		ID:      id,          // 唯一ID
		Name:    serviceName, // 服务名称
		Address: serviceIP,   // 服务地址
		Port:    servicePort, // 服务端口
		Tags:    []string{"GRPC"},
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", serviceIP, servicePort), // gRPC健康检查地址
			Interval:                       "10s",                                        // 检查间隔
			Timeout:                        "10s",
			DeregisterCriticalServiceAfter: "30s", // 检查失败后注销时间
		},
	}

	// 注册服务
	if err := client.Agent().ServiceRegister(registration); err != nil {
		log.Printf("Register fail,err:%v", err)
		panic(fmt.Sprintf("Service register failed: %v", err))
	}
	log.Println("Service registered successfully")
	return nil
}

// 获取服务实例
func DiscoverService(consulIp, serviceName string, consulPort int) (string, error) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", consulIp, consulPort)
	client, err := api.NewClient(config)
	if err != nil {
		return "", err
	}

	// 查询健康服务实例
	services, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no available services")
	}
	log.Printf("Service %s has %d services", serviceName, len(services))

	// 随机选择一个实例
	service := services[rand.Intn(len(services))]
	addr := fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)
	return addr, nil
}

// serviceDiscover in client example.
//func GetService(serviceName string) (conn *grpc.ClientConn, err error) {
//	if serviceName == "" {
//		serviceName = config.MyConfig.Grpc.Name
//	}
//	service, err := consuleServiceDiscovery.DiscoverService(config.MyConfig.Consul.Host, serviceName, config.MyConfig.Consul.Port)
//	if err != nil {
//		log.Println("discover service error", err)
//		return
//	}
//	conn, err = grpc.NewClient(service, grpc.WithInsecure())
//	if err != nil {
//		log.Fatalf("did not connect: %v", err)
//	}
//	return conn, err
//}
