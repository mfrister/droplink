package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

func cleanup() {
	if !(config.CleanupEveryMinutes > 0) {
		log.Printf("CLEANUP_EVERY_MINUTES is not set - cleanup disabled")
		return
	}
	if !(config.CleanupAfterHours > 0) {
		log.Printf("CLEANUP_AFTER_HOURS is not set - cleanup disabled")
		return
	}
	go func() {
		interval := time.Duration(config.CleanupEveryMinutes) * time.Minute
		log.Printf("Cleaning up every %v", interval)
		ticker := time.NewTicker(interval)

		for {
			<-ticker.C
			log.Printf("Cleaning up...")

			removeIfOlderThan := time.Now().Add(
				(-1 * time.Duration(config.CleanupAfterHours) * time.Hour))

			entries, err := ioutil.ReadDir(config.DataDir)
			if err != nil {
				log.Printf("Error: Cleanup failed to read dir: %v", err)
				continue
			}

			for _, entry := range entries {
				keyPath := path.Join(config.DataDir, entry.Name())
				perKeyEntries, err := ioutil.ReadDir(keyPath)
				if err != nil {
					log.Printf("Error: Cleanup failed to read dir: %v", err)
					continue
				}

				if len(perKeyEntries) != 1 {
					log.Printf("Error: Cleanup expected 1 file for key '%s', but got %d",
						entry.Name(), len(perKeyEntries))
					continue
				}

				fileEntry := perKeyEntries[0]
				if fileEntry.ModTime().Before(removeIfOlderThan) {
					filePath := path.Join(config.DataDir, entry.Name(), fileEntry.Name())

					log.Printf("  Removing: %s/%s ('%s')",
						entry.Name(), fileEntry.Name(), filePath)

					if err := os.Remove(filePath); err != nil {
						log.Printf("    File removal failed: %v", err)
						continue
					}
					if err := os.Remove(keyPath); err != nil {
						log.Printf("    Dir removal failed: %v", err)
						continue
					}
				}
			}
		}
	}()
}
