package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Favorites struct {
	UUID     string
	Name     string
	HostName string
	Icon     string
	Port     string
	Color    string
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
func saveNewList(favorites Favorites) {
	oldList := getListFromJson()
	var newFav []Favorites
	var myKey string
	for key, val := range oldList {
		newFav = append(val, favorites)
		myKey = key
	}
	newList := map[string][]Favorites{
		myKey: newFav,
	}

	rankingsJson, _ := json.Marshal(newList)
	os.WriteFile("Favorites.json", rankingsJson, 0644)
}

func getUUID(urlPath *url.URL) string {
	idSplit := strings.Split(urlPath.Path, "/")
	id := idSplit[2]
	return id
}

func removeFromList(id string) {
	oldList := getListFromJson()
	var myKey string
	var newFav []Favorites
	for key, val := range oldList {
		myKey = key
		for _, v := range val {
			if v.UUID != id {
				newFav = append(newFav, v)
			}
		}
	}
	newList := map[string][]Favorites{
		myKey: newFav,
	}

	rankingsJson, _ := json.Marshal(newList)
	os.WriteFile("Favorites.json", rankingsJson, 0644)
}

func editList(favorite Favorites) {
	oldList := getListFromJson()
	var myKey string
	var newFav []Favorites
	for key, val := range oldList {
		myKey = key
		for _, v := range val {
			if v.UUID != favorite.UUID {
				newFav = append(newFav, v)
			}
		}
		newFav = append(newFav, favorite)
	}
	newList := map[string][]Favorites{
		myKey: newFav,
	}
	rankingsJson, _ := json.Marshal(newList)
	os.WriteFile("Favorites.json", rankingsJson, 0644)
}

func getOneFav(id string) Favorites {
	all := getListFromJson()
	var targetFav Favorites
	for _, val := range all {
		for _, v := range val {
			if v.UUID == id {
				targetFav = v
			}
		}
	}
	return targetFav
}

func main() {
	fmt.Println("Running ...")

	getFavFomJson := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("www/html/index.html"))
		favorites := getListFromJson
		tmpl.Execute(w, favorites())
	}

	addFav := func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		name := r.PostFormValue("name")
		hostname := r.PostFormValue("hostname")
		icon := r.PostFormValue("icon")
		port := r.PostFormValue("port")
		color := r.PostFormValue("color")
		tmpl := template.Must(template.ParseFiles("www/html/index.html"))
		tmpl.ExecuteTemplate(w, "favorite-list-element", Favorites{UUID: id.String(), Name: name, HostName: hostname, Icon: icon, Port: port, Color: color})
		go saveNewList(Favorites{UUID: id.String(), Name: name, HostName: hostname, Icon: icon, Port: port, Color: color})
	}

	addFavForm := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("www/html/favorite-form.html"))
		tmpl.Execute(w, nil)
	}

	deleteFav := func(w http.ResponseWriter, r *http.Request) {
		id := getUUID(r.URL)
		go removeFromList(id)
	}

	editFavForm := func(w http.ResponseWriter, r *http.Request) {
		id := getUUID(r.URL)
		target := getOneFav(id)
		tmpl := template.Must(template.ParseFiles("www/html/favorite-edit.html"))
		tmpl.ExecuteTemplate(w, "favorite-edit-element", Favorites{UUID: target.UUID, Name: target.Name, HostName: target.HostName, Icon: target.Icon, Port: target.Port, Color: target.Color})
	}

	editFav := func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		hostname := r.PostFormValue("hostname")
		icon := r.PostFormValue("icon")
		port := r.PostFormValue("port")
		color := r.PostFormValue("color")
		id := getUUID(r.URL)
		target := getOneFav(id)
		tmpl := template.Must(template.ParseFiles("www/html/index.html"))
		tmpl.ExecuteTemplate(w, "favorite-list-element", Favorites{UUID: target.UUID, Name: name, HostName: hostname, Icon: icon, Port: port, Color: color})
		go editList(Favorites{UUID: target.UUID, Name: name, HostName: hostname, Icon: icon, Port: port, Color: color})
	}

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./www/css"))))
	http.HandleFunc("/", getFavFomJson)
	http.HandleFunc("/add-favorite/", addFav)
	http.HandleFunc("/edit-favorite/", editFav)
	http.HandleFunc("/delete-favorite/", deleteFav)
	http.HandleFunc("/edit-favorite-form/", editFavForm)
	http.HandleFunc("/add-favorite-form/", addFavForm)

	log.Fatal(http.ListenAndServe(":8182", nil))
}

//TODO make edit/create/delete in one function
