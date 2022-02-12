package netutils

import (
	"context"
	"io"
	"time"

	"github.com/alex-eftimie/utils"
	"github.com/fatih/color"
)

type contextKey string

var (
	// ContextKeyPipeTimeout holds timeout for piping needs
	ContextKeyPipeTimeout = contextKey("pipeTimeout")
)

// RunPiper pipes data between upstream and downstream and closes one when the other closes
func RunPiper(ctx context.Context, downstream io.ReadWriteCloser, upstream io.ReadWriteCloser) {
	var dur time.Duration
	iff := ctx.Value(ContextKeyPipeTimeout)
	if iff == nil {
		dur = time.Duration(2 * time.Hour)
	} else {
		dur = iff.(time.Duration)
	}
	// ctx := context.Background()
	// dur := time.Duration(10 * time.Second)

	done := false
	cancel := func() {
		if done {
			return
		}
		done = true
		ctx.Done()
		downstream.Close()
		upstream.Close()
		utils.Debug(999, "Closing sockets")
	}
	idleTimeoutPipe(ctx, downstream, upstream, dur, cancel)
	// go idleTimeoutPipe(ctx, upstream, downstream, dur, cancel)
	// pipe content
	// go Transfer(downstream, upstream)
	// Transfer(upstream, downstream)
}

// Transfer just copies from source to destination then closes both
func Transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
	utils.Debug(999, "Closing sockets")
}

func idleTimeoutPipe(ctx context.Context, dst io.ReadWriter, src io.ReadWriter, timeout time.Duration,
	cancel context.CancelFunc) (written int64, err error) {
	read := make(chan int)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(timeout):

				utils.Debug(999, color.RedString("idleTimeoutPipe timeout reached"))
				cancel()
				return
			case <-read:
			}
		}
	}()

	go func() {
		defer cancel()

		buf := make([]byte, 32*1024)
		for {
			nr, er := dst.Read(buf)
			if nr > 0 {
				read <- nr
				nw, ew := src.Write(buf[0:nr])
				written += int64(nw)
				if ew != nil {
					err = ew
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}
		}
	}()

	defer cancel()
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			read <- nr
			nw, ew := dst.Write(buf[0:nr])
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
