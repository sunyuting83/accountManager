package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	pathaaa, _ := GetCurrentPath()
	// fmt.Println(pathaaa)
	GetPlotsFileNumber(pathaaa)
	// fmt.Println(x)
}

func GetCurrentPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(path)
	return dir, nil
}

// GetPlotsFileNumber Get Plots File Number
func GetPlotsFileNumber(dirname string) []string {
	// Remove the trailing path separator if dirname has.
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))

	infos, err := os.ReadDir(dirname)
	if err != nil {
		return nil
	}

	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		patha := dirname + string(os.PathSeparator) + info.Name()
		if info.IsDir() {
			tmp := GetPlotsFileNumber(patha)
			if err != nil {
				return nil
			}
			paths = append(paths, tmp...)
			continue
		}
		fileSuffix := path.Ext(info.Name())
		if strings.Contains(info.Name(), ".log") && fileSuffix != ".log" {
			os.Remove(patha)
			paths = append(paths, patha)
		}
	}
	return paths
}
