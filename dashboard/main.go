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
)

type Favorites struct {
	Name     string
	HostName string
	Icon     string
	Port     string
	Color    string
}

func getList() map[string][]Favorites {
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
	oldList := getList()
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

func removeFromList(myurl *url.URL) {
	oldList := getList()
	var myKey string
	var newFav []Favorites
	id := strings.Split(myurl.Path, "/")
	for key, val := range oldList {
		myKey = key
		for _, v := range val {
			if v.Name != id[2] {
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

func main() {
	fmt.Println("Running ...")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("www/html/index.html"))
		favorites := getList
		tmpl.Execute(w, favorites())
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		hostname := r.PostFormValue("hostname")
		icon := r.PostFormValue("icon")
		port := r.PostFormValue("port")
		color := r.PostFormValue("color")
		tmpl := template.Must(template.ParseFiles("www/html/index.html"))
		tmpl.ExecuteTemplate(w, "favorite-list-element", Favorites{Name: name, HostName: hostname, Icon: icon, Port: port, Color: color})
		go saveNewList(Favorites{Name: name, HostName: hostname, Icon: icon, Port: port, Color: color})
	}

	h3 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("www/html/favorite-form.html"))
		tmpl.Execute(w, nil)
	}

	h4 := func(w http.ResponseWriter, r *http.Request) {
		go removeFromList(r.URL)
	}

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./www/css"))))
	http.HandleFunc("/", h1)
	http.HandleFunc("/add-favorite/", h2)
	http.HandleFunc("/remove-favorite/", h4)
	http.HandleFunc("/modal", h3)

	log.Fatal(http.ListenAndServe(":8182", nil))
}

//TODO generate a UUID for each fav to delete or edit it
