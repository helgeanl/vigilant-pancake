package queue

import (
	def "definitions"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// runBackup loads queue data from file if file exists once, and saves
// backups whenever its asked to. If it finds a non-empty backed up queue,
// its internal requests are added to the local queue, and its external requests
// are reassigned by sending as new requests on the network.
func runBackup(outgoingMsg chan<- def.Message) {

	const filename ="elevator_backup.dat"
	var backup queueType
	backup.loadFromDisk(filename)
	// Read last time backup was modified
	fileStat,err := os.Stat(filename);
	if  err != nil{
		log.Println(def.ColR, err, def.ColN)
	}

	// Resend all hall requests found in backup, and add cab requests to queue:
	for floor := 0; floor < def.NumFloors; floor++ {
		for btn := 0; btn < def.NumButtons; btn++ {
			if backup.hasRequest(floor, btn) {
				if btn == def.BtnCab {
					AddRequest(floor, btn, def.LocalIP)
				} else if time.Now().After(fileStat.ModTime().Add(def.RequestTimeoutDuration))  {
					outgoingMsg <- def.Message{Category: def.NewRequest, Floor: floor, Button: btn}
				}
			}
		}
	}
	go func() {
		for {
			<-takeBackup
				log.Println(def.ColG, "Take Backup", def.ColN)
			if err := queue.saveToDisk(filename); err != nil {
				log.Println(def.ColR, err, def.ColN)
			}
		}
	}()
}

// saveToDisk saves a queue to disk.
func (q *queueType) saveToDisk(filename string) error {

	data, err := json.Marshal(&q)
	log.Println(data)
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
func (q *queueType) loadFromDisk(filename string) error {
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
	printQueue()
	return nil
}
