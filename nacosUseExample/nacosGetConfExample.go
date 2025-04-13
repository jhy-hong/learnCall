package nacosUseExample

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
)

func GetNacosConfigExample(scheme, ipaddr string, port uint64, namespaceId, dataId, group string) (configStr string, err error) {
	// 1. 配置Nacos服务器参数
	serverConfigs := []constant.ServerConfig{
		{
			Scheme:      scheme,    // Nacos服务器IP
			ContextPath: "./nacos", // Nacos配置上下文
			IpAddr:      ipaddr,    // 协议（若启用HTTPS则改为https）
			Port:        port,      // NaCos 服务的上下文路径
		},
	}
	// 2. 客户端配置（含鉴权） 如果我们配置鉴权的话
	clientConfig := constant.ClientConfig{
		NamespaceId:         namespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./tmp/nacos/log",   // 日志目录
		CacheDir:            "./tmp/nacos/cache", // 缓存目录
	}
	// 3. 创建配置客户端
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ServerConfigs: serverConfigs,
			ClientConfig:  &clientConfig,
		},
	)
	if err != nil {
		log.Println("create nacos config client failed,err:", err)
		return
	}
	// 从 NaCos 获取配置内容
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId, // 配置的 DataId
		Group:  group,  // 配置的 Group
	})
	// 5. 监听配置变更（可选）
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("\n配置发生变更：")
			fmt.Println(data)
		},
	})
	if err != nil {
		log.Fatal("监听配置失败:", err)
	}
	return content, err

}
