package main

import "fmt"

func main() {
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	log.Printf("Watching for environment objects...")

	// Sync environments on startup
	err := syncEnvironments()

	if err != nil {
		log.Fatal(err)
	}

	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	// Watch for events that add, modify, or delete Certificate definitions and
	// process them asynchronously.
	log.Println("Watching for certificate events.")
	wg.Add(1)
	watchEnvironmentEvents(doneChan, &wg)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			log.Printf("Shutdown signal received, exiting...")
			close(doneChan)
			wg.Wait()
			os.Exit(0)
		}
	}
}
