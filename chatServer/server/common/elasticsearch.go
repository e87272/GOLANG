package common

import (
	"context"
	"errors"
)

func Eschatinit(IndexName string) error {

	ctx := context.Background()

	// log.Printf("EsInit IndexName : %s\n", IndexName)
	// Use the IndexExists service to check if a specified index exists.
	exists, err := Elasticclient.IndexExists(IndexName).Do(ctx)
	// log.Printf("Elasticclient exists got: %v", exists)
	if err != nil {
		return err
	}
	if exists {
		// log.Printf("%s is exists", IndexName)
		return nil
	}

	// Create a new index.
	// Create index
	createIndex, err := Elasticclient.CreateIndex(IndexName).Do(ctx)
	if err != nil {
		// log.Printf("expected CreateIndex to succeed; got: %v", err)
	}
	if createIndex == nil {
		// log.Printf("expected result to be != nil; got: %v", createIndex)
	}

	// log.Printf("createIndex response; got: %v", createIndex)

	mapping := ` {
					"properties" : {
						"historyUuid":{
							"type":"keyword"
						},
						"myUuid":{
							"type":"keyword"
						},
						"myPlatformUuid":{
							"type":"keyword"
						},
						"myPlatform":{
							"type":"keyword"
						},
						"chatTarget":{
							"type":"keyword"
						},
						"stamp":{
							"type":"date"
						},
						"message":{
							"type":"text"
						},
						"style":{
							"type":"text"
						},
						"ip":{
							"type":"keyword"
						},
						"forwardChatMessage":{
							"type":"text"
						}
					}
				}`

	putresp, err := Elasticclient.PutMapping().Index(IndexName).BodyString(mapping).Do(context.TODO())
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

	getresp, err := Elasticclient.GetMapping().Index(IndexName).Do(context.TODO())
	if err != nil {
		// log.Printf("expected get mapping to succeed; got: %v", err)
	}
	if getresp == nil {
		// log.Printf("expected get mapping response; got: %v", getresp)
	}

	// log.Printf("get mapping response; got: %v", getresp)

	_, ok := getresp[IndexName]
	if !ok {
		// log.Printf("expected JSON root to be of type map[string]interface{}; got: %#v", props)
	}

	// log.Printf("props response; got: %v", props)

	return nil
}

func Essyserrorinit(IndexName string) error {

	ctx := context.Background()

	// log.Printf("EsInit IndexName : %s\n", IndexName)
	// Use the IndexExists service to check if a specified index exists.
	exists, err := Elasticclient.IndexExists(IndexName).Do(ctx)
	// log.Printf("Elasticclient exists got: %v", exists)
	if err != nil {
		return err
	}
	if exists {
		// log.Printf("%s is exists", IndexName)
		return nil
	}

	// Create a new index.
	// Create index
	createIndex, err := Elasticclient.CreateIndex(IndexName).Do(ctx)
	if err != nil {
		// log.Printf("expected CreateIndex to succeed; got: %v", err)
	}
	if createIndex == nil {
		// log.Printf("expected result to be != nil; got: %v", createIndex)
	}

	// log.Printf("createIndex response; got: %v", createIndex)

	mapping := ` {
					"properties" : {
						"userUuid":{
							"type":"keyword"
						},
						"code":{
							"type":"keyword"
						},
						"message":{
							"type":"text"
						},
						"error":{
							"type":"text"
						},
						"stamp":{
							"type":"date"
						}
					}
				}`

	putresp, err := Elasticclient.PutMapping().Index(IndexName).BodyString(mapping).Do(context.TODO())
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

	getresp, err := Elasticclient.GetMapping().Index(IndexName).Do(context.TODO())
	if err != nil {
		// log.Printf("expected get mapping to succeed; got: %v", err)
	}
	if getresp == nil {
		// log.Printf("expected get mapping response; got: %v", getresp)
	}

	// log.Printf("get mapping response; got: %v", getresp)

	_, ok := getresp[IndexName]
	if !ok {
		// log.Printf("expected JSON root to be of type map[string]interface{}; got: %#v", props)
	}

	// log.Printf("props response; got: %v", props)

	return nil
}
func Esdelete(IndexName string) error {

	ctx := context.Background()

	deleteIndex, err := Elasticclient.DeleteIndex(IndexName).Do(ctx)
	if err != nil {
		return err
	}
	if !deleteIndex.Acknowledged {
		return errors.New("delete index not acknowledged")
	}
	return nil
}

func Esinsert(IndexName string, data string) error {

	//加鎖 加鎖 加鎖
	// log.Printf("Mutexelastic :ElasticclientLock\n")
	Mutexelastic.Lock()
	defer func() {
		// log.Printf("Mutexelastic :ElasticclientUNLock\n")
		Mutexelastic.Unlock() // 完成後記得 解鎖 解鎖 解鎖
	}()
	_, err := Elasticclient.Index().
		Index(IndexName).
		BodyString(data).
		Do(context.Background())
	if err != nil {
		// log.Printf("Elasticclient err : %+v\n", err)
		return err
	}

	// log.Printf("Elasticclient res : %+v\n", res)
	return nil
}
