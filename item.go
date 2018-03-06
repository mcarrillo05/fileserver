package fileserver

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
)

type typeFile int

const (
	dirType typeFile = iota
	fileType
)

var (
	sizes = [...]string{"Bytes", "KB", "MB", "GB", "TB"}
)

type itemsResponse struct {
	Root  Item   `json:"root"`
	Items []Item `json:"items"`
}

//Item is a file or directory.
type Item struct {
	Name       string   `json:"name"`
	Size       int64    `json:"size"`
	SizeString string   `json:"size_string"`
	Date       string   `json:"date"`
	Type       typeFile `json:"type"`
}

func getItems(root string, countFiles bool) (items []Item) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		i := Item{
			Name: info.Name(),
		}
		if !info.IsDir() {
			i.Type = fileType
			i.Size = info.Size()
			i.SizeString = formatSize(info.Size())
			i.Date = info.ModTime().Format("2006-01-02 03:04PM")
		} else {
			i.Name += "/"
		}
		items = append(items, i)
		if info.IsDir() && path != root {
			return filepath.SkipDir
		}
		return nil
	})
	for i := range items {
		if i == 0 {
			if items[i].Type == dirType {
				items[i].Size = int64(len(items) - 1)
				items[i].SizeString = strconv.FormatInt(items[0].Size, 10)
			}
		} else {
			if countFiles {
				items[i].countFiles(filepath.Join(root, items[i].Name))
			}
		}
	}
	return items
}

func (i *Item) countFiles(root string) {
	if i.Type == fileType {
		return
	}
	i.Size = 0
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			i.Size++
		}
		return nil
	})
	i.SizeString = strconv.FormatInt(i.Size, 10)
}

func formatSize(sz int64) string {
	if sz == 0 {
		return "0 Bytes"
	}
	idx := int(math.Log(float64(sz)) / math.Log(1024))
	if idx == 0 {
		return strconv.FormatInt(sz, 10) + " " + sizes[idx]
	}
	return strconv.FormatFloat(float64(sz)/math.Pow(1024, float64(idx)), 'f', 2, 64) + " " + sizes[idx]
}
