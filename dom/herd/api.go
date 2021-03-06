package herd

import (
	"github.com/Symantec/Dominator/dom/mdb"
	"github.com/Symantec/Dominator/lib/image"
	"github.com/Symantec/Dominator/sub/scanner"
	"io"
	"log"
	"sync"
	"time"
)

const (
	statusUnknown = iota
	statusConnecting
	statusFailedToConnect
	statusWaitingToPoll
	statusPolling
	statusFailedToPoll
	statusSubNotReady
	statusImageNotReady
	statusFetching
	statusFailedToFetch
	statusWaitingForNextPoll
	statusUpdating
	statusFailedToUpdate
	statusSynced
)

type HtmlWriter interface {
	WriteHtml(writer io.Writer)
}

type Sub struct {
	herd                         *Herd
	hostname                     string
	requiredImage                string
	plannedImage                 string
	busyMutex                    sync.Mutex
	busy                         bool
	fileSystem                   *scanner.FileSystem
	generationCount              uint64
	generationCountAtChangeStart uint64
	status                       uint
	lastPollStartTime            time.Time
	lastPollSucceededTime        time.Time
	lastShortPollDuration        time.Duration
	lastFullPollDuration         time.Duration
}

type Herd struct {
	sync.RWMutex            // Protect map and slice mutations.
	imageServerAddress      string
	logger                  *log.Logger
	htmlWriters             []HtmlWriter
	nextSubToPoll           uint
	subsByName              map[string]*Sub
	subsByIndex             []*Sub // Sorted by Sub.hostname.
	imagesByName            map[string]*image.Image
	makeConnectionSemaphore chan bool
	pollSemaphore           chan bool
	currentScanStartTime    time.Time
	previousScanDuration    time.Duration
}

func NewHerd(imageServerAddress string, logger *log.Logger) *Herd {
	return newHerd(imageServerAddress, logger)
}

func (herd *Herd) MdbUpdate(mdb *mdb.Mdb) {
	herd.mdbUpdate(mdb)
}

func (herd *Herd) PollNextSub() bool {
	return herd.pollNextSub()
}

func (herd *Herd) StartServer(portNum uint, daemon bool) error {
	return herd.startServer(portNum, daemon)
}

func (herd *Herd) AddHtmlWriter(htmlWriter HtmlWriter) {
	herd.addHtmlWriter(htmlWriter)
}
