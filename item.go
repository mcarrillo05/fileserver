package fileserver

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
)

//TypeItem defines if item is directory or file.
type TypeItem int

const (
	//DirType is a directory.
	DirType TypeItem = iota
	//FileType is a file.
	FileType
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
	ID         string   `json:"ID,omitempty"`
	Name       string   `json:"name"`
	Size       int64    `json:"size"`
	SizeString string   `json:"size_string"`
	Date       string   `json:"date"`
	Type       TypeItem `json:"type"`
}

//GetItems returns all files and directories, if countFiles is true, process will count all files of each directory at first level.
func GetItems(root string, countFiles bool) (items []Item) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		i := Item{
			Name: info.Name(),
		}
		if !info.IsDir() {
			i.Type = FileType
			i.Size = info.Size()
			i.SizeString = formatSize(info.Size())
			i.Date = info.ModTime().Format("2006-01-02 03:04PM")
		} else {
			i.Name += "/"
			if countFiles {
				i.countFiles(filepath.Join(root, i.Name))
			}
		}
		items = append(items, i)
		if info.IsDir() && path != root {
			return filepath.SkipDir
		}
		return nil
	})
	if len(items) > 0 && items[0].Type == DirType {
		items[0].Size = int64(len(items) - 1)
		items[0].SizeString = strconv.FormatInt(items[0].Size, 10)
	}
	return items
}

func (i *Item) countFiles(root string) {
	if i.Type == FileType {
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
