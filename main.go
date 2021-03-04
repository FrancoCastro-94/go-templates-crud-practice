package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

type player struct {
	Id       int
	Name     string
	LastName string
	Number   string
}

type players []player

var allPlayers = players{
	{
		Id:       1,
		Name:     "Franco",
		LastName: "Catro",
		Number:   "34",
	},
}

var tmpl = template.Must(template.ParseGlob("views/*"))

//Index view
func Index(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "Index", allPlayers)

}

//Delete by player id
func searchById(id int) player {
	for _, p := range allPlayers {
		if p.Id == id {
			return p
		}
	}
	return player{}
}

//update by player id
func updateById(id int, name, lastName, number string) {
	playerDelete := deleteById(id)
	newP := player{}
	newP.Id = id
	newP.Name = name
	newP.LastName = lastName
	newP.Number = number
	allPlayers = append(playerDelete, newP)
}

//deleteById return a slice without player deleted
func deleteById(id int) players {
	result := []player{}
	for _, p := range allPlayers {
		if p.Id != id {
			result = append(result, p)
		}
	}
	return result
}

//add player to slice allPlayers
func addPlayer(id int, name, lastName, number string) {
	p := player{}
	p.Id = id
	p.Name = name
	p.LastName = lastName
	p.Number = number
	allPlayers = append(allPlayers, p)
	log.Println(allPlayers)
}

//View player
func View(w http.ResponseWriter, r *http.Request) {
	nID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal()
	}
	p := searchById(nID)
	tmpl.ExecuteTemplate(w, "Show", p)
}

//New view
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

//Insert player
func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := len(allPlayers) + 2
		name := r.FormValue("name")
		lastName := r.FormValue("lastName")
		number := r.FormValue("number")
		addPlayer(id, name, lastName, number)
	}
	http.Redirect(w, r, "/", 301)
}

//Edit Player
func Edit(w http.ResponseWriter, r *http.Request) {
	nID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	p := searchById(nID)
	tmpl.ExecuteTemplate(w, "Edit", p)
}

//Update player
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nID, _ := strconv.Atoi(r.URL.Query().Get("id"))
		name := r.FormValue("name")
		lastName := r.FormValue("lastName")
		number := r.FormValue("number")
		updateById(nID, name, lastName, number)

		log.Println("UPDATE: Name: " + name + " | Last Name: " + lastName + " | Number: " + number)
	}
	http.Redirect(w, r, "/", 301)
}

//Delete player
func Delete(w http.ResponseWriter, r *http.Request) {
	nID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println("algo acurrio")
	}
	allPlayers = deleteById(nID)
	log.Println("Delete")
	http.Redirect(w, r, "/", 301)
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("Server started on port: " + port)

	http.HandleFunc("/", Index)
	http.HandleFunc("/show", View)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/new", New)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/delete", Delete)

	http.ListenAndServe(":"+port, nil)

}
