package main

import (
	"log"
	"sync"
)

func syncEnvironments() error {
	environments, err := getEnvironments()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, env := range environments {
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

func watchEnvironmentEvents(done chan struct{}, wg *sync.WaitGroup) {
	events, watchErrs := monitorEnvironmentEvents()
	go func() {
		for {
			select {
			case event := <-events:
				err := processEnvironmentEvent(event)
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

func processEnvironmentEvent(c EnvironmentEvent) error {
	switch {
	case c.Type == "ADDED":
		return processEnvironment(c.Object)
	case c.Type == "DELETED":
		return deleteEnvironment(c.Object)
	}
	return nil
}

func deleteEnvironment(e Environment) error {
	log.Println("Deleting Environment")
	return nil
}
func processEnvironment(e Environment) error {
	log.Println("Processing Environment")
	return nil
}
