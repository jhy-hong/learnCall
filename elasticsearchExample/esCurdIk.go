package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strings"
	"time"
)

const (
	ElasticsearchURL = "https://localhost:9200"
	IndexName        = "go-8-test"
	Username         = "elastic"              // 默认用户名
	Password         = "your_password"        // 运行 ./bin/elasticsearch-reset-password -u elastic 获取
	CACert           = "/path/to/http_ca.crt" // 从ES配置目录获取
)

type Document struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}

func main() {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"}, // 替换为您的 Elasticsearch 地址
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	document := map[string]interface{}{
		"user":    "john_doe",
		"message": "Hello, Elasticsearch!",
	}

	jsonData, err := json.Marshal(document)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}
	createIndex(es)
	indexDocument(es, "test_index", "texto1", string(jsonData))

	query := `{
		"query": {
			"match": {
				"user": "Elasticsearch"
			}
		}
	}`
	searchDocuments(es, "doe", query)
}

// 创建文档
func indexDocument(es *elasticsearch.Client, indexName, docID, jsonData string) {
	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       strings.NewReader(jsonData),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%s", res.Status(), "your_document_id")
	} else {
		log.Printf("[%s] Document ID=%s indexed successfully", res.Status(), "your_document_id")
	}
}

// 创建索引并设置 IK 分词器
func createIndex(es *elasticsearch.Client) {
	indexName := "test_index"
	body := `{
		"settings": {
			"analysis": {
				"tokenizer": {
					"ik_max_word": {
						"type": "ik_max_word"
					}
				},
				"analyzer": {
					"default": {
						"type": "custom",
						"tokenizer": "ik_max_word"
					}
				}
			}
		}
	}`
	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(body),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}
	fmt.Println(res)
	defer res.Body.Close()
}

/*
updateData := map[string]interface{}{
    "doc": map[string]interface{}{
        "message": "Updated message content",
    },
}

jsonData, err := json.Marshal(updateData)
if err != nil {
    log.Fatalf("Error marshaling update data: %s", err)
}


复杂更新
scriptData := map[string]interface{}{
    "script": map[string]interface{}{
        "source": "ctx._source.counter += params.count",
        "params": map[string]interface{}{
            "count": 4,
        },
    },
}

jsonData, err := json.Marshal(scriptData)
if err != nil {
    log.Fatalf("Error marshaling script data: %s", err)
}

*/

// 更新文档
func updateDocument(es *elasticsearch.Client, indexName, docID, jsonData string) {
	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       strings.NewReader(jsonData),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error updating document: %s", err)
	}
	defer res.Body.Close()
}

// 删除文档
func deleteDocument(es *elasticsearch.Client, indexName, docID string) {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: docID,
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error deleting document: %s", err)
	}
	defer res.Body.Close()
}

// 查询文档
func searchDocuments(es *elasticsearch.Client, indexName, query string) {
	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  strings.NewReader(query),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error searching documents: %s", err)
	}
	defer res.Body.Close()
	// 处理查询结果
	fmt.Println(res)
}
