package elasticSearch

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/olivere/elastic"
)

var elasticClient *elastic.Client
var mutexElastic = new(sync.Mutex)

func EsInit() {

	var err error
	// log.Printf("elasticSearchHost : %v\n", os.Getenv("elasticSearchHost"))
	elasticClient, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(os.Getenv("elasticSearchHost")))

	for err != nil {
		log.Printf("Elasticsearch err %+v\n", err)
		timer := time.NewTimer(time.Second)
		select {
		case <-timer.C:
			elasticClient, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(os.Getenv("elasticSearchHost")))
		}
	}

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := elasticClient.ElasticsearchVersion(os.Getenv("elasticSearchHost"))
	if err != nil {
		log.Printf("Elasticsearch err %+v\n", err)
	}
	log.Printf("Elasticsearch version %s\n", esversion)

	sysLogInit()
	sysErrorLogInit()

}

func EsDelete(IndexName string) error {

	ctx := context.Background()

	deleteIndex, err := elasticClient.DeleteIndex(IndexName).Do(ctx)
	if err != nil {
		return err
	}
	if !deleteIndex.Acknowledged {
		return errors.New("delete index not acknowledged")
	}
	return nil
}

func EsInsert(IndexName string, data string) error {

	//加鎖 加鎖 加鎖
	// log.Printf("mutexElastic :elasticClientLock\n")
	mutexElastic.Lock()
	defer func() {
		// log.Printf("mutexElastic :elasticClientUNLock\n")
		mutexElastic.Unlock() // 完成後記得 解鎖 解鎖 解鎖
	}()
	_, err := elasticClient.Index().
		Index(IndexName).
		BodyString(data).
		Do(context.Background())
	if err != nil {
		// log.Printf("elasticClient err : %+v\n", err)
		return err
	}

	// log.Printf("elasticClient res : %+v\n", res)
	return nil
}
