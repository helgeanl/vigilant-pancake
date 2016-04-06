package queue

import (
	def "config"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// runBackup loads queue data from file if file exists once, and saves
// backups whenever its asked to. If it finds a non-empty backed up queue,
// its internal order are added to the local queue, and its external orders
// are reassigned by sending as new orders on the network.
func runBackup(outgoingMsg chan<- def.Message) {
	const filename = "lift_backup"

	var backup queue
	backup.loadFromDisk(filename)

	// Resend all orders found on loaded backup file:
	if !backup.isEmpty() {
		for f := 0; f < def.NumFloors; f++ {
			for b := 0; b < def.NumButtons; b++ {
				if backup.isOrder(f, b) {
					if b == def.BtnInside {
						AddLocalOrder(f, b)
					} else {
						outgoingMsg <- def.Message{Category: def.NewOrder, Floor: f, Button: b}
					}
				}
			}
		}
	}
	go func() {
		for {
			<-takeBackup
			if err := local.saveToDisk(filename); err != nil {
				log.Println(def.ColR, err, def.ColN)
			}
		}
	}()
}

// saveToDisk saves a queue to disk.
func (q *queue) saveToDisk(filename string) error {

	data, err := json.Marshal(&q)
	if err != nil {
		log.Println(def.ColR, "json.Marshal() error: Failed to backup.", def.ColN)
		return err
	}
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		log.Println(def.ColR, "ioutil.WriteFile() error: Failed to backup.", def.ColN)
		return err
	}
	return nil
}

// loadFromDisk checks if a file of the given name is available on disk, and
// saves its contents to a queue if the file is present.
func (q *queue) loadFromDisk(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		log.Println(def.ColG, "Backup file found, processing...", def.ColN)

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println(def.ColR, "loadFromDisk() error: Failed to read file.", def.ColN)
		}
		if err := json.Unmarshal(data, q); err != nil {
			log.Println(def.ColR, "loadFromDisk() error: Failed to Unmarshal.", def.ColN)
		}
	}
	return nil
}
