package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Command struct {
	docs []document
	idx  index
}

func indexHandler(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "search.html")
}

func (command Command) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	search := request.FormValue("search")
	start := time.Now()
	matchedIDs := command.idx.search(search)
	//log.Printf("Найдено %d документов за %v", len(matchedIDs), time.Since(start))

	type Record struct {
		Id   string
		Text string
	}

	type Message struct {
		Header  string
		Records []Record
	}

	message := Message{
		Header:  "По запросу `" + search + "` найдено " + fmt.Sprint(len(matchedIDs)) + " документов за " + fmt.Sprint(time.Since(start)),
		Records: []Record{},
	}

	for _, id := range matchedIDs {
		doc := command.docs[id]
		record := Record{Id: "Документ №" + fmt.Sprint(id), Text: doc.Text}
		message.Records = append(message.Records, record)
		//log.Printf("%d\t%s\n", id, doc.Text)
	}
	tmpl, _ := template.ParseFiles("find.html")
	tmpl.Execute(response, message)
}

func main() {

	dumpPath := "ruwiki-latest-abstract.xml.gz"
	log.Println("Старт")
	start := time.Now()
	docs, err := loadDocuments(dumpPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Загружено из архива %d документов за %v", len(docs), time.Since(start))

	start = time.Now()
	idx := make(index)
	idx.add(docs)
	log.Printf("Проиндексировано %d документов за %v", len(docs), time.Since(start))

	command := Command{docs: docs, idx: idx}
	http.HandleFunc("/", indexHandler)
	http.Handle("/search", command)
	http.ListenAndServe(":8000", nil)
}
