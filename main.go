package main

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type Item struct {
	Added             string         `xml:"Added"`
	LockData          bool           `xml:"LockData"`
	LocalTitle        string         `xml:"LocalTitle"`
	RunningTime       int            `xml:"RunningTime"`
	Genres            []Genre        `xml:"Genres"`
	PlaylistItems     []PlaylistItem `xml:"PlaylistItems>PlaylistItem"`
	Shares            []Share        `xml:"Shares>Share"`
	PlaylistMediaType string         `xml:"PlaylistMediaType"`
}

type Genre struct {
	Genre string `xml:"Genre"`
}

type PlaylistItem struct {
	XMLName   xml.Name `xml:"PlaylistItem"`
	Path      string   `xml:"Path"`
	DateAdded string   `xml:"DateAdded,attr"`
}

type Path struct {
	Path      string `xml:"Path"`
	DateAdded string `xml:"DateAdded,attr"`
}

type Share struct {
	UserId  string `xml:"UserId"`
	CanEdit bool   `xml:"CanEdit"`
}

var smartPlaylist Item
var nameOfPlaylist = "DWCFPlaylist.xml"
var playlistTitle = "Dirty Workz Copyright Free - Gen"
var userId = "01c3cfad3190498bb91cc9d4080608d9"
var canEdit = true
var folderToWatch = "/home/austin/Files/Jellyfin/Music/Dirty Workz - Copyright Free"
var subFolder = "/Jellyfin/Music/Dirty Workz - Copyright Free"
var indexToTakeOnwards = 4

// TODO Populate this data into a hash table, that way I can get a linear look up
var acceptableAudioTypes = []string{".mp3", ".wav", ".m4a", ".flac", ".mp4", ".wma", ".ogg", ".aac"}
// var acceptableAudioTypes = []string{".mp3"}

//TODO: Create a playlist with just arguments that I provide, and that way I can create them on the fly
//-t title
//-d directory
//-r recursive

func main() {
	files, err := getAllMusicFilesFromFolder(folderToWatch)
	if err != nil {
		fmt.Println(err)
	}
	// readXmlFile()

	for _, file := range files {
		splitString := strings.Split(file, "/")
		finalPath := fmt.Sprintf("/%s",strings.Join(splitString[indexToTakeOnwards:], "/"))

    fmt.Println(isAllowedFileExtension(filepath.Ext(file)))
    if isAllowedFileExtension(filepath.Ext(file)) {
      smartPlaylist.PlaylistItems = append(smartPlaylist.PlaylistItems, PlaylistItem{Path: finalPath, DateAdded: time.Now().Format("2006.01.02 15:04:05")})
      fmt.Println("Added file to playlist")
    }
	}

	smartPlaylist.Added = time.Now().Format("2006.01.02 15:04:05")
	smartPlaylist.LockData = false
	smartPlaylist.LocalTitle = playlistTitle
	//TODO Calculate running time of the playlist
	smartPlaylist.RunningTime = 0
	smartPlaylist.PlaylistMediaType = "Audio"
	//TODO Figure out Genres added to the playlist
	smartPlaylist.Genres = []Genre{{Genre: "Hardstyle"}}
	smartPlaylist.Shares = []Share{{UserId: userId, CanEdit: true}}
  writeXML()
}

func getAllMusicFilesFromFolder(folder string) ([]string, error) {
	var files []string
	//get all the files from only the folder
	entries, err := os.ReadDir(folder)
	if err != nil {
		return []string{}, err
	}

	for _, entry := range entries {
		if isAllowedFileExtension(filepath.Ext(entry.Name())) {
			files = append(files, fmt.Sprintf("%s/%s", folder, entry.Name()))
		}
	}

	return files, nil
}

func addFileToPlaylist() {

}

func createSmartPlaylist() {
	readConfig()
	// fmt.Println(nameOfPlaylist)
	// fmt.Println(playlistTitle)
	// fmt.Println(userId)
	// fmt.Println(canEdit)
	// fmt.Println(folderToWatch)
	// fmt.Println(subFolder)
	// fmt.Println(indexToTakeOnwards)

	readXmlFile()
	smartPlaylist.Added = time.Now().Format("2006.01.02 15:04:05")
	smartPlaylist.LockData = false
	smartPlaylist.LocalTitle = playlistTitle
	//TODO Calculate running time of the playlist
	smartPlaylist.RunningTime = 0
	smartPlaylist.PlaylistMediaType = "Audio"
	//TODO Figure out Genres added to the playlist
	smartPlaylist.Genres = []Genre{{Genre: "todo"}}
	smartPlaylist.Shares = []Share{{UserId: userId, CanEdit: true}}

	folders, err := getAllFoldersToWatch()
	if err != nil {
		log.Fatal(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, folder := range folders {
		if folder[0] != '.' {
			err = watcher.Add(folder)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
					if event.Has(fsnotify.Create) {
						ext := filepath.Ext(event.Name)
						if isAllowedFileExtension(ext) == false {
							continue
						}

						splitString := strings.Split(event.Name, "/")
						finalPath := fmt.Sprintf("%s/%s", subFolder, strings.Join(splitString[indexToTakeOnwards:], "/"))

						//I might have to use a lock or a channel with this.
						if doesTrackAlreadyExist(finalPath) {
							continue
						}

						smartPlaylist.PlaylistItems = append(smartPlaylist.PlaylistItems, PlaylistItem{Path: finalPath, DateAdded: time.Now().Format("2006.01.02 15:04:05")})
						writeXML()
						fmt.Println("New file added to xml")
					}
				}

			case err := <-watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()

	<-make(chan struct{})

}

func readConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	file, err := os.ReadFile(fmt.Sprintf("%s/%s", homeDir, ".config/jellyfin-smart-playlist/config.txt"))
	if err != nil {
		fmt.Println("No config file found, please create a config file at ~/.config/jellyfin-smart-playlist/config.txt")
		panic(err)
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if strings.Contains(line, "name_of_playlist") {
			nameOfPlaylist = strings.Split(line, "=")[1]
		} else if strings.Contains(line, "playlist_title") {
			playlistTitle = strings.Split(line, "=")[1]
		} else if strings.Contains(line, "user_id") {
			userId = strings.Split(line, "=")[1]
		} else if strings.Contains(line, "can_edit") {
			canEdit = strings.Split(line, "=")[1] == "true"
		} else if strings.Contains(line, "folder_to_watch") {
			folderToWatch = strings.Split(line, "=")[1]
		} else if strings.Contains(line, "sub_folder") {
			subFolder = strings.Split(line, "=")[1]
		} else if strings.Contains(line, "index_to_take_onwards") {
			indexToTakeOnwards, err = strconv.Atoi(strings.Split(line, "=")[1])
			if err != nil {
				panic(err)
			}
		}
	}
}

func readXmlFile() {
	file, err := os.ReadFile(nameOfPlaylist)
	if err != nil {
		fmt.Println("No file found, creating new one")
		smartPlaylist = Item{}
		writeXML()
		return
	}
	xml.Unmarshal(file, &smartPlaylist)
}

func getAllFoldersToWatch() ([]string, error) {
	var folders []string
	err := filepath.Walk(folderToWatch, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			folders = append(folders, path)
		}
		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return folders, nil
}

func isAllowedFileExtension(ext string) bool {
	for _, allowedExt := range acceptableAudioTypes {
		if strings.EqualFold(ext, allowedExt) {
			return true
		}
	}
	return false
}

func doesTrackAlreadyExist(path string) bool {
	for _, track := range smartPlaylist.PlaylistItems {
		if track.Path == path {
			return true
		}
	}

	return false
}

func sortPlaylist() {
	//TODO Sort the playlist, so that's more of a queue and not just an array
}

func writeXML() {
	file, _ := xml.MarshalIndent(smartPlaylist, "", " ")
	file = []byte(xml.Header + string(file))
	_ = os.WriteFile(nameOfPlaylist, file, 0644)

}
