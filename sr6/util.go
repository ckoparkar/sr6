package sr6

import (
	"bytes"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// OverwriteFile overwrites file at *path* with *content*
func OverwriteFile(path, content string) error {
	os.Remove(path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(content)
	return nil
}

// CorrectHostname ensures that server hostname ends with suffix,
// and adds random chars before the suffix if we are setting it now.
func CorrectHostname(suffix string) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(hostname, suffix) {
		// change hostname.
		hostname = RandString(5) + "." + suffix
		if err := SetHostname(hostname); err != nil {
			return "", err
		}
	}
	return hostname, nil
}

// SetHostname sets hostname to *name*
func SetHostname(name string) error {
	cmd := exec.Command("hostname", name)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Printf("[INFO] Setting hostname to %s\n", name)
	return nil
}

// RandString returns random string *n* chars long
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}
