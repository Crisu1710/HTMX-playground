package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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

func main() {
	fmt.Println("Running ...")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		favorites := getList
		tmpl.Execute(w, favorites())
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		hostname := r.PostFormValue("hostname")
		icon := r.PostFormValue("icon")
		port := r.PostFormValue("port")
		color := r.PostFormValue("color")
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "favorite-list-element", Favorites{Name: name, HostName: hostname, Icon: icon, Port: port, Color: color})
		saveNewList(Favorites{Name: name, HostName: hostname, Icon: icon, Port: port, Color: color})
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-favorite/", h2)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
