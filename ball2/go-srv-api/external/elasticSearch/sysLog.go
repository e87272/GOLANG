package elasticSearch

import (
	"context"
	"os"
)

func sysLogInit() error {

	indexName := os.Getenv("sysLog")
	ctx := context.Background()

	// log.Printf("EsInit indexName : %s\n", indexName)
	// Use the IndexExists service to check if a specified index exists.
	exists, err := elasticClient.IndexExists(indexName).Do(ctx)
	// log.Printf("elasticClient exists got: %v", exists)
	if err != nil {
		return err
	}
	if exists {
		// log.Printf("%s is exists", indexName)
		return nil
	}

	// Create a new index.
	// Create index
	createIndex, err := elasticClient.CreateIndex(indexName).Do(ctx)
	if err != nil {
		// log.Printf("expected CreateIndex to succeed; got: %v", err)
	}
	if createIndex == nil {
		// log.Printf("expected result to be != nil; got: %v", createIndex)
	}

	// log.Printf("createIndex response; got: %v", createIndex)

	mapping := ` {
					"properties" : {
						"apiName":{
							"type":"keyword"
						},
						"msg":{
							"type":"text"
						},
						"stamp":{
							"type":"date"
						}
					}
				}`

	putresp, err := elasticClient.PutMapping().Index(indexName).BodyString(mapping).Do(context.TODO())
	if err != nil {
		// log.Printf("expected put mapping to succeed; got: %v", err)
	}
	if putresp == nil {
		// log.Printf("expected put mapping response; got: %v", putresp)
	}
	if !putresp.Acknowledged {
		// log.Printf("expected put mapping ack; got: %v", putresp.Acknowledged)
	}

	// log.Printf("putresp response; got: %v", putresp)

	getresp, err := elasticClient.GetMapping().Index(indexName).Do(context.TODO())
	if err != nil {
		// log.Printf("expected get mapping to succeed; got: %v", err)
	}
	if getresp == nil {
		// log.Printf("expected get mapping response; got: %v", getresp)
	}

	// log.Printf("get mapping response; got: %v", getresp)

	_, ok := getresp[indexName]
	if !ok {
		// log.Printf("expected JSON root to be of type map[string]interface{}; got: %#v", props)
	}

	// log.Printf("props response; got: %v", props)

	return nil
}
