package main

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Favorites struct {
	UUID     string
	Name     string
	Icon     string
	Protocol string
	HostName string
	Port     string
	Path     string
	Note     string
	Color    string
}

func init() {
	log.Println("start init ...")
	checkFile, err := os.Stat("Favorites.json")
	if err != nil {
		log.Println(err)
	}
	if errors.Is(err, os.ErrNotExist) || checkFile.Size() == 0 {
		log.Println("File is empty. Create init file ...")
		newList := map[string][]Favorites{
			"Favorites": nil,
		}
		rankingsJson, _ := json.Marshal(newList)
		os.WriteFile("Favorites.json", rankingsJson, 0644)
	} else if checkFile.Size() != 0 {
		log.Println("Json file found. Starting ...")
	}
}

func getUUID(urlPath *url.URL) string {
	idSplit := strings.Split(urlPath.Path, "/")
	id := idSplit[len(idSplit)-1]
	return id
}

func getListFromJson() map[string][]Favorites {
	content, err := os.ReadFile("Favorites.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var favorites map[string][]Favorites
	err = json.Unmarshal(content, &favorites)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return favorites
}

func genNewJsonList(favorite Favorites, method string, id string) {
	jsonList := getListFromJson()
	var newFavList []Favorites
	var myKey string
	for key, valFromJson := range jsonList {
		myKey = key
		if method == "DELETE" {
			for _, v := range valFromJson {
				if v.UUID != id {
					newFavList = append(newFavList, v)
				}
			}
		} else if method == "PUT" {
			for _, v := range valFromJson {
				if v.UUID != favorite.UUID {
					newFavList = append(newFavList, v)
				}
			}
			newFavList = append(newFavList, favorite)
		} else if method == "POST" {
			newFavList = append(valFromJson, favorite)
		} else {
			log.Println(method)
			log.Fatal("unknown method")
		}
	}
	newList := map[string][]Favorites{
		myKey: newFavList,
	}

	rankingsJson, _ := json.Marshal(newList)
	os.WriteFile("Favorites.json", rankingsJson, 0644)
}

func getOneFavByID(id string) Favorites {
	jsonList := getListFromJson()
	var targetFav Favorites
	for _, valFromJson := range jsonList {
		for _, v := range valFromJson {
			if v.UUID == id {
				targetFav = v
			}
		}
	}
	return targetFav
}

func main() {
	log.Println("Running ...")

	getFavFomJson := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("www/html/index.gohtml"))
		favorites := getListFromJson
		tmpl.Execute(w, favorites())
	}

	addFav := func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		name := r.PostFormValue("name")
		protocol := r.PostFormValue("protocol")
		path := r.PostFormValue("path")
		note := r.PostFormValue("note")
		hostname := r.PostFormValue("hostname")
		icon := r.PostFormValue("icon")
		port := r.PostFormValue("port")
		color := r.PostFormValue("color")
		tmpl := template.Must(template.ParseFiles("www/html/index.gohtml"))
		tmpl.ExecuteTemplate(w, "favorite-list-element", Favorites{UUID: id.String(), Name: name, Icon: icon, Protocol: protocol, Path: path, HostName: hostname, Port: port, Note: note, Color: color})
		go genNewJsonList(Favorites{UUID: id.String(), Name: name, Icon: icon, Protocol: protocol, Path: path, HostName: hostname, Port: port, Note: note, Color: color}, r.Method, "")
	}

	addFavForm := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("www/html/favorite-form.gohtml"))
		tmpl.ExecuteTemplate(w, "favorite-edit-element", Favorites{})
	}

	deleteFav := func(w http.ResponseWriter, r *http.Request) {
		id := getUUID(r.URL)
		go genNewJsonList(Favorites{}, r.Method, id)
	}

	editFavForm := func(w http.ResponseWriter, r *http.Request) {
		id := getUUID(r.URL)
		target := getOneFavByID(id)
		tmpl := template.Must(template.ParseFiles("www/html/favorite-form.gohtml"))
		tmpl.ExecuteTemplate(w, "favorite-edit-element", Favorites{UUID: target.UUID, Name: target.Name, Protocol: target.Protocol, Path: target.Path, HostName: target.HostName, Icon: target.Icon, Port: target.Port, Note: target.Note, Color: target.Color})
	}

	editFav := func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		protocol := r.PostFormValue("protocol")
		path := r.PostFormValue("path")
		note := r.PostFormValue("note")
		hostname := r.PostFormValue("hostname")
		icon := r.PostFormValue("icon")
		port := r.PostFormValue("port")
		color := r.PostFormValue("color")
		id := getUUID(r.URL)
		target := getOneFavByID(id)
		tmpl := template.Must(template.ParseFiles("www/html/index.gohtml"))
		tmpl.ExecuteTemplate(w, "favorite-list-element", Favorites{UUID: target.UUID, Name: name, Icon: icon, Protocol: protocol, Path: path, HostName: hostname, Port: port, Note: note, Color: color})
		go genNewJsonList(Favorites{UUID: target.UUID, Name: name, Icon: icon, Protocol: protocol, Path: path, HostName: hostname, Port: port, Note: note, Color: color}, r.Method, "")
	}

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./www/css"))))
	http.HandleFunc("/", getFavFomJson)
	http.HandleFunc("/form/favorite/add", addFavForm)
	http.HandleFunc("/form/favorite/edit/", editFavForm)
	http.HandleFunc("/favorite/add/", addFav)
	http.HandleFunc("/favorite/edit/", editFav)
	http.HandleFunc("/favorite/delete/", deleteFav)

	log.Fatal(http.ListenAndServe(":8182", nil))
}
