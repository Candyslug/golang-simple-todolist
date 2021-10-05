package main

// Simple todo list api (to learn golang)
// Code not perfect.. nor great (of course)

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

type Item struct {
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Id          int    `json:"id"`
	Priority    bool   `json:"priority"`
}

type Index struct {
	I int `json:"i"`
}

type IndexToggle struct {
	I          int  `json:"i"`
	Completion bool `json:"completion"`
}

var Data = []Item{}

// Hardcoded list for testing
func initItems() {
	Data = append(Data, Item{
		Description: "Feed the dog",
		Completed:   false,
		Priority:    true,
	})
	Data = append(Data, Item{
		Description: "Eat some candy",
		Completed:   true,
		Priority:    false,
	})
	Data = append(Data, Item{
		Description: "Do some coding!",
		Completed:   false,
		Priority:    true,
	})
}

// =========================== //

func showItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request method"))
		return
	}

	j, err := json.Marshal(Data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func showItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request method"))
		return
	}
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid request"))
		return
	}

	b, id, err := strIsDigit(parts[len(parts)-1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	if b != true {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request, index NaN"))
		return
	}
	if int(id) > int((len(Data) - 1)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Invalid request, index [%v] out of bounds [%v]", id, (len(Data) - 1))))
		return
	}
	j, err := json.Marshal(Data[id])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request method"))
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("Invalid request, unsupported media type"))
		return
	}
	var buf Item
	err = json.Unmarshal(b, &buf)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Printf("prev Data n : %v", len(Data))
	Data = append(Data, buf)
	fmt.Printf("New Data n : %v", len(Data))
}

func delItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request method"))
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("Invalid request, unsupported media type"))
		return
	}
	var buf Index
	err = json.Unmarshal(b, &buf)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if buf.I > len(Data)-1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request, index out of bounds"))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Printf("Removed [%v] from items\n", buf.I)
	Data = rmItemFromSlice(Data, buf.I)
}

func toggleCompletion(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request method"))
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("Invalid request, unsupported media type"))
		return
	}
	var buf IndexToggle
	err = json.Unmarshal(b, &buf)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if buf.I > len(Data)-1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request, index out of bounds"))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Printf("Removed [%v] from items\n", buf.I)
	fmt.Printf("Toggled [%v] completion status\n", buf.I)
	fmt.Printf("------ %v -> %v\n", Data[buf.I].Completed, buf.Completion)
	Data[buf.I].Completed = buf.Completion
}

func main() {

	fmt.Println("Server Started")

	initItems()

	//get
	http.HandleFunc("/items", showItems)
	http.HandleFunc("/items/", showItem)
	//post
	http.HandleFunc("/items/add", addItem)
	http.HandleFunc("/items/del", delItem)           // json data contains the index
	http.HandleFunc("/items/mark", toggleCompletion) // json data contains the index (not url)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func strIsDigit(s string) (bool, int, error) {
	for _, i := range s {
		if !unicode.IsDigit(i) {
			return false, -1, nil
		}
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return false, -1, err
	}
	return true, i, nil
}

func rmItemFromSlice(s []Item, i int) []Item {
	return append(s[:i], s[i+1:]...)
}
