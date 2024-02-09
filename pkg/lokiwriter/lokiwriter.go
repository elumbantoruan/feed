package lokiwriter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/loki/pkg/push"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// LokiWriter writes the log to Loki.
// It maybe used for the system which doesn't or cannot utilize a promtail
// such as Cronjob
type LokiWriter struct {
	grpcEndpoint string
	labels       string
	conn         *grpc.ClientConn
	entries      chan push.Entry
	buffer       chan push.Entry
	ticker       *time.Ticker
	wait         time.Duration
	batchSize    int
}

type Labels map[string]string

func NewLokiWriter(grpcEndpoint string, labels Labels) (writer *LokiWriter, shutdown func(), err error) {
	conn, err := grpc.Dial(grpcEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}
	lbls := "{"
	for k, v := range labels {
		lbls += fmt.Sprintf(`%s = "%s"`, k, v)
		lbls += ","
	}
	if len(lbls) > 1 {
		lbls = strings.TrimRight(lbls, ",")
	}
	lbls += "}"

	wait := time.Duration(500 * time.Millisecond)

	lw := &LokiWriter{
		grpcEndpoint: grpcEndpoint,
		labels:       lbls,
		conn:         conn,
		entries:      make(chan push.Entry),
		buffer:       make(chan push.Entry, 100),
		wait:         wait,
		ticker:       time.NewTicker(wait),
		batchSize:    2,
	}

	// run LokiWriter in the background waiting for log entries to be send into Loki
	go lw.run()

	return lw,
		func() {
			close(lw.entries) // close entries channel
			lw.exporter(true) // drain the buffer for log entry
			conn.Close()      // close loki grpc connection
		},
		nil
}

func (lw LokiWriter) run() {

	for {
		select {
		case entry, ok := <-lw.entries:
			if !ok {
				continue
			}
			lw.buffer <- entry

			for len(lw.buffer) >= lw.batchSize {
				lw.ticker.Stop()
				lw.exporter(false)
			}
			lw.ticker.Reset(lw.wait)

		case <-lw.ticker.C:
			// ticker elapsed and push whatever in the logs
			// process the log that's in the buffer
			for len(lw.buffer) > 0 {
				lw.exporter(false)
			}
		}
	}
}

func (lw LokiWriter) exporter(drain bool) {
	var entries []push.Entry
	var count = 0
	if len(lw.buffer) < lw.batchSize || drain {
		count = len(lw.buffer)
	} else {
		count = lw.batchSize
	}

	for i := 0; i < count; i++ {
		entries = append(entries, <-lw.buffer)
	}

	err := lw.Push(context.Background(), entries)
	if err != nil {
		fmt.Println(err)
	}
}

// Write writes the log.  Writes implement io.Writer
// that are used in log writer such as slog
func (lw LokiWriter) Write(p []byte) (n int, err error) {
	line := "stdout F " + string(p)
	lw.entries <- push.Entry{
		Timestamp: time.Now(),
		Line:      line,
	}
	return len(line), nil
}

// Push pushes the message to Loki.
// It follows the term used in Loki
func (lw LokiWriter) Push(ctx context.Context, entries []push.Entry) error {
	req := push.PushRequest{
		Streams: []push.Stream{
			{
				Labels:  lw.labels,
				Entries: entries,
			},
		},
	}

	pc := push.NewPusherClient(lw.conn)
	_, err := pc.Push(ctx, &req)
	if err != nil {
		return err
	}
	return nil
}
