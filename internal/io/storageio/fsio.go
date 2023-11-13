package storageio

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gnames/bhlquest/pkg/ent/storage"
	"github.com/gnames/gnsys"
)

func (s *storageio) Pages(itemID uint) ([]storage.Page, error) {
	var res []storage.Page
	sID := fmt.Sprintf("%06d", itemID)
	if len(sID) != 6 {
		return res, fmt.Errorf("ItemID %d is not correct", itemID)
	}
	path := filepath.Join(s.cfg.BHLDir, sID[0:3], sID)
	exists, empty, err := gnsys.DirExists(path)
	if err != nil || !exists || empty {
		err = fmt.Errorf("path '%s' does not exist or is empty: %w", path, err)
		return nil, err
	}

	return s.pathToPages(path, itemID)
}

func (s *storageio) pathToPages(
	path string,
	itemID uint,
) ([]storage.Page, error) {
	es, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var start, end uint
	var res []storage.Page
	for i := range es {
		name := es[i].Name()
		if filepath.Ext(name) != ".txt" {
			fmt.Println("BAD")
			continue
		}
		pageID, err := nameToPageID(es[i].Name())
		if err != nil {
			return nil, err
		}
		filePath := filepath.Join(path, name)
		bs, err := os.ReadFile(filePath)
		txt := string(bs)
		end = start + uint(len(txt))
		page := storage.Page{
			ID:       pageID,
			ItemID:   itemID,
			FileName: name,
			Text:     txt,
			Start:    start,
			End:      end,
		}
		start = end
		res = append(res, page)
	}
	return res, nil
}

func nameToPageID(fName string) (uint, error) {
	es := strings.Split(fName, "-")
	if len(es) != 3 {
		return 0, fmt.Errorf("bad file name %s", fName)
	}
	num, err := strconv.Atoi(es[1])
	if err != nil {
		return 0, err
	}

	return uint(num), nil
}
