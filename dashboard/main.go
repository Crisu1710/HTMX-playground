package main

import (
	"encoding/json"
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
		go genNewJsonList(Favorites{UUID: id.String(), Name: name, HostName: hostname, Icon: icon, Port: port, Color: color}, r.Method, "")
	}

	addFavForm := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("www/html/favorite-add.html"))
		tmpl.Execute(w, nil)
	}

	deleteFav := func(w http.ResponseWriter, r *http.Request) {
		id := getUUID(r.URL)
		go genNewJsonList(Favorites{}, r.Method, id)
	}

	editFavForm := func(w http.ResponseWriter, r *http.Request) {
		id := getUUID(r.URL)
		target := getOneFavByID(id)
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
		target := getOneFavByID(id)
		tmpl := template.Must(template.ParseFiles("www/html/index.html"))
		tmpl.ExecuteTemplate(w, "favorite-list-element", Favorites{UUID: target.UUID, Name: name, HostName: hostname, Icon: icon, Port: port, Color: color})
		go genNewJsonList(Favorites{UUID: target.UUID, Name: name, HostName: hostname, Icon: icon, Port: port, Color: color}, r.Method, "")
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

//TODO make edit/create/delete in one function
