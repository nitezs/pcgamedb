package utils

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func MegaDownload(url string, path string) (string, []string, error) {
	stat, err := os.Stat("torrent")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("torrent", 0755)
			if err != nil {
				return "", nil, err
			}
		} else {
			return "", nil, err
		}
	}
	if !stat.IsDir() {
		os.Remove("torrent")
		err = os.Mkdir("torrent", 0755)
		if err != nil {
			return "", nil, err
		}
	}
	cmd := exec.Command("mega-get", url, path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", nil, err
	}
	pathRegex := regexp.MustCompile(`(?i)Download finished: (.*)`)
	pathRegexRes := pathRegex.FindAllStringSubmatch(out.String(), -1)
	if len(pathRegexRes) == 0 {
		return "", nil, errors.New("Mega download failed")
	}
	pathRegexRes[0][1] = strings.TrimSpace(pathRegexRes[0][1])
	res, err := walkDir(pathRegexRes[0][1])
	if err != nil {
		return "", nil, err
	}
	return pathRegexRes[0][1], res, nil
}

func walkDir(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, file := range files {
		if file.IsDir() {
			subFiles, err := walkDir(filepath.Join(path, file.Name()))
			if err != nil {
				return nil, err
			}
			res = append(res, subFiles...)
		} else {
			res = append(res, filepath.Join(path, file.Name()))
		}
	}
	return res, nil
}
