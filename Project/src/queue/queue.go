package queue

import (
	def "definitions"
	"fmt"
	"time"
)

type requestStatus struct {
	status bool
	addr   string       `json:"-"`
	timer  *timer.Timer `json:"-"`
}

type queue struct {
	matrix [def.Numfloors][def.NumButtons]requestStatus
}

//make a request inactive
const inactive = requestStatus{status: false, addr: "", timer: nil}

func (q *queue) hasRequest(floor, btn int) bool {
	return q.matrix[floor][btn].status
}
/// -------------------

var queue queue

# Setting PATH for Python 3.3
# The orginal version is saved in .bash_profile.pysave
PATH="/Library/Frameworks/Python.framework/Versions/3.3/bin:${PATH}"
export PATH

export XUGGLE_HOME="/usr/local/xuggler"

export TRACKER_HOME="/usr/local/tracker"

#export PS1="\[\033[36m\]\u\[\033[m\]\h:\[\033[33;1m\]\w\$ "
#export PS1='\[\e[1;32m\][\u@\h \W]\$\[\e[0m\]'
#export PS1="\[\033[36m\]\u\[\033[m\]:\[\033[m\]\$"
export PS1="\[\033[36m\]\u\[\033[m\]@\[\033[32m\]\h:\[\033[33;1m\]\w\[\033[m\]\$ "
#export PS1="\[\033[1;32m\]\u@\h\[\033[0m\]:\[\033[1;34m\]\w\[\033[0m\]# "


export CLICOLOR=1
export LSCOLORS= GxFxCxDxBxegedabagaced
alias ls='ls -GFh'
alias showFiles='defaults write com.apple.finder AppleShowAllFiles YES; killall Finder /System/Library/CoreServices/Finder.app'
alias hideFiles='defaults write com.apple.finder AppleShowAllFiles NO; killall Finder /System/Library/CoreServices/Finder.app'
export PATH=/usr/local/gnat/bin:$PATH

# Size of terminal history file
HISTFILESIZE=100000

# Colors for cat
alias ccat='pygmentize -g'
alias ccatl='pygmentize -g -O style=colorful,linenos=1'

func Init(newRequestTemp chan bool, outgoingMsg chan def.Message) {
	newRequest = newRequestTemp /// ??????
	go updateLocalQueue()
	runBackup(outgoingMsg)
	log.Println(def.ColG, "Queue initialised.", def.ColN)
}

func (q *queue) setRequest(floor, btn int, request requestStatus){
	q.matrix[floor][btn] = request.status
	// sync lights
	// take backup
	// print
}

func AddRequestAt(floor int, btn int, addr string){
	if !queue.hasRequest(floor,btn){
		queue.setRequest(floor,btn,requestStatus{floor,btn,addr,nil})
		queue.startTimer(floor, btn)
	}
}

func (q * queue) startTimer(floor, btn int){
	q.matrix[floor][btn].timer = time.NewTimer(def.RequestTimeoutDuration)
	<-q.matrix[floor][btn].timer.C
	// Wait until timeout
	RequestTimeoutChan <- def.Btnpress{floor, btn}
}

func (q * queue) stopTimer(floor, btn int){
	if q.matrix[floor][btn].timer != nil{
		q.matrix[floor][btn].timer.Stop()
	}
}

// Go through queue, and resend requests belonging to dead elevator
func ReassignRequest(addr string){

}

func RemoveOrderAt(floor int){

}

// -------------------

// requests_above
func (q *queue) hasRequestAbove(floor int) bool {
	for f := floor + 1; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if q.hasRequest(f, b) {
				return true
			}
		}
	}
	return false
}

// requests_below
func (q *queue) hasRequestsBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if q.hasRequest(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) chooseDirection(floor, dir int) int {
	switch dir {
	case def.DirUp:
		if q.hasRequestsAbove(floor){
			return def.DirUp
		} else if q.hasRequestsBelow(floor){
			return def.DirDown
		} else{
			return def.DirStop
		}
	case def.DirDown, def.DirStop:
		if q.hasRequestsBelow(floor) {
			return def.DirDown
		} else if q.hasRequestsAbove(floor) {
			return def.DirUp
		} else {
			return def.DirStop
		}
	default:
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Printf("%sChooseDirection(): called with invalid direction %d, returning stop%s\n", def.ColR, dir, def.ColN)
		return 0
	}
}

func (q *queue) shouldStop(floor, dir int) bool {
	switch dir {
	case def.DirDown:
		return
			q.hasRequest(floor, def.BtnHallDown) ||
			q.hasRequest(floor, def.BtnCab) ||
			!q.hasRequestsBelow(floor)
	case def.DirUp:
		return
			q.hasRequest(floor, def.BtnHallUp) ||
			q.hasRequest(floor, def.BtnCab) ||
			!q.hasRequestsAbove(floor)
	case def.DirStop:
	default:
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Fatalln(def.ColR, "This direction doesn't exist", def.ColN)
	}
	return false
}
