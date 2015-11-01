package sr6

import "os"

func overwriteFile(path, content string) error {
	os.Remove(path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(content)
	return nil
}
