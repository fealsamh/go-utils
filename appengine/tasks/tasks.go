package tasks

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/taskqueue"
)

// Add ...
func Add(ctx context.Context, path string, payload []byte, delay time.Duration) (*taskqueue.Task, error) {
	if appengine.IsAppEngine() {
		return taskqueue.Add(ctx, &taskqueue.Task{
			Path:    path,
			Method:  http.MethodPost,
			Payload: payload,
			Delay:   delay,
		}, "")
	}
	go func() {
		time.Sleep(delay)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(payload))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		req.Header.Set("X-Appengine-TaskName", "task")
		var cl http.Client
		resp, err := cl.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		resp.Body.Close()
	}()
	return nil, nil
}
