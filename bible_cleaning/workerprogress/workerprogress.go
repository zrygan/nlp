package workerprogress

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"github.com/schollz/progressbar/v3"
)

type QueenConfig struct {
	IsDetailed bool // if true, print detailed per-worker progress
	ReportInterval int // milliseconds
	UseProgressBar bool // if true, use progress bar for global progress
}

type QueenContext struct {
	DoneWorkers    *atomic.Int32
	TotalWorkers   int
	StartTime      time.Time
	Quit           chan struct{}
	Wg             *sync.WaitGroup
	ProgressCh     chan WorkerProgressMsg // send-only channel
	Config				 *QueenConfig
}

type WorkerProgressMsg struct {
	WorkerID string  // e.g. "en-fr" for parallel language worker
	Percent  float32 // 0.0 → 1.0
	Status   string  // optional: text message

}

type WorkerProgressContext struct {
	Wg        *sync.WaitGroup
	Progress  chan<- WorkerProgressMsg // send-only channel
	WorkerID  string
	TotalWork int // optional: total units of work
}



func NewQueenContext(totalWorkers int, config *QueenConfig) *QueenContext {
	q := &QueenContext{
		DoneWorkers:    &atomic.Int32{},
		TotalWorkers:   totalWorkers,
		StartTime:      time.Now(),
		Quit:           make(chan struct{}),
		Wg:             &sync.WaitGroup{},
		ProgressCh:     make(chan WorkerProgressMsg, 100), // buffered channel
		Config:				 config,
	}

	return q;
}

func (q *QueenContext) CreateWorkerContext(workerID string) WorkerProgressContext {
	q.Wg.Add(1)
	fmt.Printf("Launching worker %s...\n", workerID)

	workerCtx := WorkerProgressContext{
			Wg:       q.Wg,
			Progress: q.ProgressCh,
			WorkerID: workerID,
	}

	return workerCtx;
}

func (q *QueenContext) RunReporter() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(q.Config.ReportInterval))
	defer ticker.Stop()
	defer fmt.Println("\n[Reporter] Stopped.")
	// defer close(q.Quit)

	// track per-worker state
	progressMap := make(map[string]WorkerProgressMsg)

	bar := progressbar.NewOptions(q.TotalWorkers,
		progressbar.OptionSetDescription("Processing..."),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "#", SaucerPadding: "-", BarStart: "[", BarEnd: "]"}),
		
	)

	for {
		select {
		case <-ticker.C:
			done := q.DoneWorkers.Load()
			pct := float64(done) / float64(q.TotalWorkers) * 100
			elapsed := time.Since(q.StartTime).Truncate(time.Second)
			fmt.Printf("\rGlobal: %d/%d (%.1f%%) elapsed %s", done, q.TotalWorkers, pct, elapsed)

		case msg, ok := <-q.ProgressCh:
			if !ok {
				fmt.Println("\n[Reporter] Progress channel closed.")
				return
			}

			if q.Config.IsDetailed {
				progressMap[msg.WorkerID] = msg
				fmt.Printf("\n[Reporter] Worker %s: %.1f%% - %s\n", msg.WorkerID, msg.Percent*100, msg.Status)
			}


			if msg.Percent >= 1.0 {
				newDone := q.DoneWorkers.Add(1)
				bar.Add(1)
				fmt.Println("")
				if newDone >= int32(q.TotalWorkers) {
					fmt.Println("\n[Reporter] All workers done.")
					return
				}
			}

		case <-q.Quit:
			return
		}
	}
}
