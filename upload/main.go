package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/docker/go-units"
	"github.com/go-chi/chi/v5"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func genRandomName(prefix string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return prefix + string(b)
}

func genName64(name string) string {
	name64 := make([]byte, base64.StdEncoding.EncodedLen(len(name)))
	base64.StdEncoding.Encode(name64, []byte(name))
	return string(name64)
}

func main() {
	r := chi.NewRouter()

	logger := log.New(os.Stdout, "DEBUG: ", log.Lshortfile)

	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		f, h, err := r.FormFile("upload")
		if err != nil {
			logger.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// make a pipe
		pr, pw := io.Pipe()
		// create dir
		dir := "upload"
		err = os.Mkdir(dir, 0755)
		if err != nil {
			logger.Println(err)
			return
		}
		fname := genName64(h.Filename)
		file, err := os.Create(dir + "/" + fname)
		if err != nil {
			logger.Println(err)
			return
		}
		go func() {
			defer pw.Close()
			// 4 Kib buffer
			b := make([]byte, 4*units.KiB)
			buf := bytes.NewBuffer(b)
			for {
				// read from http socket
				if _, err = io.ReadAtLeast(f, buf.Bytes(), 1); err != nil {
					// no bytes were read
					if err == io.ErrUnexpectedEOF || err == io.EOF {
						break
					}
					logger.Println(err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				// write to a pipe
				if _, err = pw.Write(buf.Bytes()); err != nil {
					logger.Println(err)
					break
				}
			}
		}()
		b := make([]byte, 4*units.KiB)
		buf := bytes.NewBuffer(b)
		for {
			if _, err = pr.Read(buf.Bytes()); err != nil {
				if err == io.EOF {
					break
				}
				logger.Println(err)
				return
			}
			_, err := file.Write(buf.Bytes())
			if err != nil {
				logger.Println(err)
				return
			}
		}
	})
	http.ListenAndServe(":8080", r)
}
