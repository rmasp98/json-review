package ui

import (
	"fmt"
	"log"
	"os"

	"github.com/awesome-gocui/gocui"
)

// Save saves content to file even if file already exists
func Save(filename string, content string) error {
	file, errOpen := os.Create(filename)
	defer file.Close()
	if errOpen != nil {
		return errOpen
	}
	if _, errWrite := file.Write([]byte(content)); errWrite != nil {
		return errWrite
	}
	return nil
}

// SaveUI stuff
func SaveUI(g *gocui.Gui, v *gocui.View) error {
	go saveProcess()
	return nil
}

func saveProcess() {
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

	if err := Save(filename, "Test"); err == nil {
		cui.CreatePopup("Save Successful", saveType+" data has been successfully saved to "+filename, NewConfirmPopupEditor(nil), true, false, true)
	} else {
		cui.CreatePopup("Save Failed", err.Error(), NewConfirmPopupEditor(nil), true, false, true)
	}
}

func getSaveType() (string, error) {
	var ch = make(chan string)
	content := "Choose what to save:\nRaw"
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
