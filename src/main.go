package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Finds struct {
	docs []document
	idx  index
}

var archPath string
var finds Finds

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/load", loadHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/find", findHandler)
	http.ListenAndServe(":8000", nil)
}

func indexHandler(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/index.html")
}

func loadHandler(response http.ResponseWriter, request *http.Request) {
	type Message struct {
		Loaded  string
		Indexed string
	}
	archPath = request.FormValue("archPath")
	log.Println("Старт")
	start1 := time.Now()
	docs, err := loadDocuments(archPath)
	if err != nil {
		log.Fatal(err)
	}

	start2 := time.Now()
	idx := make(index)
	idx.add(docs)
	finds.docs = docs
	finds.idx = idx
	tmpl, _ := template.ParseFiles("templates/load.html")
	tmpl.Execute(response, Message{
		Loaded:  fmt.Sprintf("Загружено из архива %d документов за %v", len(docs), time.Since(start1)),
		Indexed: fmt.Sprintf("Проиндексировано %d документов за %v", len(docs), time.Since(start2)),
	})
}

func searchHandler(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/search.html")
}

func findHandler(response http.ResponseWriter, request *http.Request) {
	phrase := request.FormValue("phrase")
	start := time.Now()
	matchedIDs := finds.idx.search(phrase)

	type Record struct {
		Id   string
		Text string
	}

	type Message struct {
		Header  string
		Records []Record
	}

	message := Message{
		Header:  fmt.Sprintf("По запросу <%s> найдено %d документов за %v", phrase, len(matchedIDs), time.Since(start)),
		Records: []Record{},
	}

	for _, id := range matchedIDs {
		doc := finds.docs[id]
		record := Record{Id: "Документ №" + fmt.Sprint(id), Text: doc.Text}
		message.Records = append(message.Records, record)
	}
	tmpl, _ := template.ParseFiles("templates/find.html")
	tmpl.Execute(response, message)
}
