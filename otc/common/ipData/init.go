package ipData

func Init() {

	ok = QueryWhiteListDB()

	for !ok {
		timer := time.NewTimer(time.Second)
		select {
		case <-timer.C:
			ok = QueryWhiteListDB()
		}
	}

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			QueryWhiteListDB()
			log.Printf("Query DB tick")
		}
	}()

}