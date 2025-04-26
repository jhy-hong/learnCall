package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	// 1. 初始化客户端
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://192.168.134.128:9200"}, // 替换为您的 Elasticsearch 地址
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 2. 构造请求体
	body := map[string]interface{}{
		"analyzer": "ik_max_word",
		"text":     "我是中国人",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("JSON 序列化失败：%s", err)
	}

	// 3. 调用 Analyze API
	res, err := es.Indices.Analyze(
		//es.Indices.Analyze.WithIndex("test_index"),              // 可选：指定索引
		es.Indices.Analyze.WithBody(bytes.NewReader(bodyBytes)), // 必须：用 WithBody 包装你的 io.Reader
		es.Indices.Analyze.WithContext(context.Background()),    // 可选：传入上下文
	)
	if err != nil {
		log.Fatalf("调用 Analyze API 失败：%s", err)
	}
	defer res.Body.Close()

	// 4. 解析结果
	var result struct {
		Tokens []struct {
			Token string `json:"token"`
		} `json:"tokens"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("解析响应失败：%s", err)
	}

	fmt.Println("分词结果：")
	for _, t := range result.Tokens {
		fmt.Println(" -", t.Token)
	}
}
