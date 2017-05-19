package main

func syncEnvironments() error {
	environments, err := getEnvironments()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, cert := range environments {
		wg.Add(1)
		go func(env Environment) {
			defer wg.Done()
			err := processEnvironment(env)
			if err != nil {
				log.Println(err)
			}
		}(env)
	}
	wg.Wait()
	return nil
}

func watchCertificateEvents(db *bolt.DB, done chan struct{}, wg *sync.WaitGroup) {
	events, watchErrs := monitorCertificateEvents()
	go func() {
		for {
			select {
			case event := <-events:
				err := processCertificateEvent(event, db)
				if err != nil {
					log.Println(err)
				}
			case err := <-watchErrs:
				log.Println(err)
			case <-done:
				wg.Done()
				log.Println("Stopped certificate event watcher.")
				return
			}
		}
	}()
}
