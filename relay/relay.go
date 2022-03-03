/*
 * @Date: 2022.03.02 10:16
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 10:16
 */

package relay

import (
	"io"
	"log"
)

var verbose bool

func logf(f string, v ...interface{}) {
	if verbose {
		log.Printf(f, v...)
	}
}

func SetVerboseMode(v bool) {
	verbose = v
}

// relay copies between left and right bidirectionally. Returns number of
// bytes copied from right to left, from left to right, and any error occurred.
func relay(left, right io.ReadWriter) (int64, int64, error) {
	type res struct {
		N   int64
		Err error
	}
	ch := make(chan res)

	go func() {
		n, err := copyHalfClose(right, left)
		ch <- res{n, err}
	}()

	n, err := copyHalfClose(left, right)
	rs := <-ch

	if err == nil {
		err = rs.Err
	}
	return n, rs.N, err
}

type closeWriter interface {
	CloseWrite() error
}

type closeReader interface {
	CloseRead() error
}

// copyHalfClose copies to dst from src and optionally closes dst for writing and src for reading.
func copyHalfClose(dst io.Writer, src io.Reader) (int64, error) {
	defer func() {
		// half-close to wake up other goroutines blocking on dst and src

		if c, ok := dst.(closeWriter); ok {
			c.CloseWrite()
		}

		if c, ok := src.(closeReader); ok {
			c.CloseRead()
		}
	}()

	return io.Copy(dst, src) // will use io.ReaderFrom or io.WriterTo shortcut if possible
}
