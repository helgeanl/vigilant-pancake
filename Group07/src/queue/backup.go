package queue

import (
	def "definitions"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

// RunBackup loads backup on startup, and saves queue whenever
// there is anything on the takeBackup channel
func RunBackup(outgoingMsg chan<- def.Message) {

	const filename = "elevator_backup.dat"
	var backup QueueType
	backup.loadFromDisk(filename)
	printQueue()

	// Read last time backup was modified
	fileStat, _ := os.Stat(filename)

	// Resend all hall requests found in backup, and add cab requests to queue:
	for floor := 0; floor < def.NumFloors; floor++ {
		for btn := 0; btn < def.NumButtons; btn++ {
			if backup.hasRequest(floor, btn) {
				if btn == def.BtnCab {
					AddRequest(floor, btn, def.LocalIP)
				// Check if time since last backup is less than RequestTimeoutDuration
				}else if !time.Now().After(fileStat.ModTime().Add(def.RequestTimeoutDuration)){
					AddRequest(floor, btn, def.LocalIP)
				}
			}
		}
	}
	go func() {
		for {
			<-takeBackup
			log.Println(def.ColG, "Take Backup", def.ColN)
			queue.saveToDisk(filename)
		}
	}()
}

// saveToDisk saves a QueueType to disk.
func (q *QueueType) saveToDisk(filename string) {
	data, _ := json.Marshal(&q)
	ioutil.WriteFile(filename, data, 0644)
}

// loadFromDisk checks if a file of the given name is available on disk, and
// saves its contents to a QueueType
func (q *QueueType) loadFromDisk(filename string) {
	if _, err := os.Stat(filename); err == nil {
		log.Println(def.ColG, "Backup file found, processing...", def.ColN)
		data, _ := ioutil.ReadFile(filename)
		json.Unmarshal(data, q)
	}
}

func printQueue() {
	fmt.Println(def.ColB, "\n*****************************")
	fmt.Println("*       Up     Down    Cab   ")
	for f := def.NumFloors - 1; f >= 0; f-- {
		s := "* " + strconv.Itoa(f+1) + "  "
		for b := 0; b < def.NumButtons; b++ {
			if queue.hasRequest(f, b) && b != def.BtnCab {
				s += "( " + queue.Matrix[f][b].Addr[12:15] + " ) "
			} else if queue.hasRequest(f, b) {
				s += "(  x  ) "
			} else {
				s += "(     ) "
			}
		}
		fmt.Println(s)
	}
	fmt.Println("*****************************\n", def.ColN)
}
