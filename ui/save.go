package ui

import (
	"fmt"
	"kube-review/nodelist"
	"kube-review/search"
	"log"
	"os"

	"github.com/awesome-gocui/gocui"
)

// SaveUI sets up processes for extracting save information from user
type SaveUI struct {
	nodeList  *nodelist.NodeList
	queryList *search.QueryList
}

// NewSaveUI stuff
func NewSaveUI(nodeList *nodelist.NodeList, queryList *search.QueryList) SaveUI {
	return SaveUI{nodeList, queryList}
}

// Save stuff
func (s SaveUI) Save(g *gocui.Gui, v *gocui.View) error {
	go s.saveProcess()
	return nil
}

func (s SaveUI) saveProcess() {
	saveType, errType := getSaveType()
	if errType != nil {
		log.Println(errType.Error())
		return
	}

	filename, errName := getFilename(saveType)
	if errName != nil {
		log.Println(errName.Error())
		return
	}

	var err error
	if saveType == "Raw" {
		// err = s.nodeList.Save(filename)
	} else if saveType == "Query" {
		err = s.queryList.Save(filename)
	} else {
		err = fmt.Errorf(saveType + " is not a valid save option")
	}

	if err == nil {
		cui.CreatePopup("Save Successful", saveType+" data has been successfully saved to "+filename, NewConfirmPopupEditor(nil), true, false, true)
	} else {
		cui.CreatePopup("Save Failed", err.Error(), NewConfirmPopupEditor(nil), true, false, true)
	}
}

func getSaveType() (string, error) {
	var ch = make(chan string)
	content := "Choose what to save:\nRaw\nQuery"
	if err := cui.CreatePopup("Save", content, NewSelectPopupEditor(ch), false, true, true); err != nil {
		return "", fmt.Errorf("Could not create popup")
	}
	saveType := <-ch
	cui.ClosePopup()
	return saveType, nil
}

func getFilename(saveType string) (string, error) {
	var ch = make(chan string)
	if err := cui.CreatePopup("Save "+saveType, "Provide filename:\n", NewWritePopupEditor(ch), true, false, true); err != nil {
		return "", fmt.Errorf("Could not create popup")
	}
	filename := <-ch
	cui.ClosePopup()
	if _, err := os.Stat(filename); err == nil {
		cui.CreatePopup("File Exists", "This file already exists. Do you want to overwrite: (Y/N)", NewConfirmPopupEditor(ch), true, false, true)
		overwrite := <-ch
		cui.ClosePopup()
		if overwrite == "n" {
			return "", fmt.Errorf("User chose not to overwrite file")
		}
	}
	return filename, nil
}
