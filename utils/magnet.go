package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/anacrolix/torrent/metainfo"
)

func ConvertTorrentToMagnet(torrent []byte) (string, string, error) {
	minfo, err := metainfo.Load(bytes.NewReader(torrent))
	if err != nil {
		return "", "", err
	}
	info, err := minfo.UnmarshalInfo()
	if err != nil {
		return "", "", err
	}
	var size int64 = info.Length
	if size == 0 {
		for _, file := range info.Files {
			size += file.Length
		}
	}
	infoHash := minfo.HashInfoBytes()
	magnet := minfo.Magnet(&infoHash, &info)
	return magnet.String(), FormatSize(size), nil
}

func FormatSize(size int64) string {
	const (
		_        = iota
		KB int64 = 1 << (10 * iota)
		MB
		GB
		TB
	)
	switch {
	case size >= GB:
		return fmt.Sprintf("%.1f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.1f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.1f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d Bytes", size)
	}
}

func SubSizeStrings(sizes []string) (string, error) {
	size := int64(0)
	for _, sizeStr := range sizes {
		sizeStr := strings.ToLower(sizeStr)
		if strings.Contains(sizeStr, "gb") {
			sizeStr = strings.ReplaceAll(sizeStr, "gb", "")
			sizeStr = strings.TrimSpace(sizeStr)
			addSize, err := strconv.ParseFloat(sizeStr, 64)
			if err != nil {
				return "", err
			}
			size += int64(addSize * 1024 * 1024 * 1024)
		}
		if strings.Contains(sizeStr, "mb") {
			sizeStr = strings.ReplaceAll(sizeStr, "mb", "")
			sizeStr = strings.TrimSpace(sizeStr)
			addSize, err := strconv.ParseFloat(sizeStr, 64)
			if err != nil {
				return "", err
			}
			size += int64(addSize * 1024 * 1024)
		}
		if strings.Contains(sizeStr, "kb") {
			sizeStr = strings.ReplaceAll(sizeStr, "kb", "")
			sizeStr = strings.TrimSpace(sizeStr)
			addSize, err := strconv.ParseFloat(sizeStr, 64)
			if err != nil {
				return "", err
			}
			size += int64(addSize * 1024)
		}
	}
	return FormatSize(size), nil
}
