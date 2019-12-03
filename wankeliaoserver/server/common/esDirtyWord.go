package common

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/olivere/elastic"
)

func Esdirtywordhistoryinit(IndexName string) error {

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
						"roomUuid":{
							"type":"keyword"
						},
						"toUuid":{
							"type":"text"
						},
						"fromUuid":{
							"type":"text"
						},
						"stamp":{
							"type":"date"
						},
						"message":{
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

func Essidetextdirtywordhistoryinit(IndexName string) error {

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
func loaddirtyword(IndexName string) {

	// 可写方式打开文件
	// file, err := os.OpenFile(
	// 	"data/dirtyword.txt",
	// 	os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
	// 	0666,
	// )

	// log.Printf("loaddirtyword")

	dirtywordData, err := ioutil.ReadFile("data/dirtyword.json")
	if err != nil {
		log.Fatal("找不到dirtyword.json")
	}

	dirtywordDataObj := [][]interface{}{}

	err = json.Unmarshal(dirtywordData, &dirtywordDataObj)
	if err != nil {
		// log.Printf("err %+v \n", err)
	}

	for _, v := range dirtywordDataObj {
		// log.Printf("%v : %+v\n", k, v[0])

		dirtyword := make(map[string]interface{})
		dirtyword["dirtyword"] = v[0]

		// 写字节到文件中
		// byteSlice := []byte(dirtyword["dirtyword"].(string) + "\n")
		// file.Write(byteSlice)

		// Index a second tweet (by string)
		dirtywordJson, _ := json.Marshal(dirtyword)

		Esinsert("dirtyword", string(dirtywordJson[:]))
	}

}

func Esdirtywordsearch(message string) (bool, error) {

	var re = regexp.MustCompile(`[^A-Za-z0-9\p{Han}]`)
	message = re.ReplaceAllString(message, ` `)

	// log.Printf("message : %s\n", message)
	//加鎖 加鎖 加鎖
	// log.Printf("Mutexelastic :ElasticclientLock\n")
	Mutexelastic.Lock()
	defer func() {
		// log.Printf("Mutexelastic :ElasticclientUNLock\n")
		Mutexelastic.Unlock() // 完成後記得 解鎖 解鎖 解鎖
	}()

	var IndexName = "dirtyword"
	// log.Printf("IndexName : %+v\n", IndexName)
	boolQ := elastic.NewBoolQuery()
	boolQ.Must(elastic.NewMatchQuery("dirtyword", message))

	// Search with a term query
	searchResult, err := Elasticclient.Search(IndexName).Query(boolQ).Do(context.Background()) // execute
	if err != nil {
		// log.Printf("searchResult err got: %v", err)
		return false, err
	}
	// log.Printf("searchResult : %+v\n", searchResult)
	// Here's how you iterate through results with full control over each step.
	// log.Printf("searchResult count: %+v\n", searchResult.Hits.TotalHits.Value)
	if searchResult.Hits.TotalHits.Value > 0 {

		// log.Printf("searchResult.Hits : %+v\n", searchResult.Hits)

		// log.Printf("maxScore : %+v\n", *searchResult.Hits.MaxScore)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index
			// log.Printf("hit : %+v\n", hit)
			var word interface{}
			err := json.Unmarshal(hit.Source, &word)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			// log.Printf("word : %+v\n", word)
		}
		return true, err
	}

	return false, nil
}
