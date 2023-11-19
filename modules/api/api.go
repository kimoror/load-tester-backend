package api

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"io"
	"load-tester-backend/modules/loader"
	"load-tester-backend/modules/model"
	"log"
	"net/http"
)

// https://cd06f14b-ba79-45b9-b89b-3006c76704d8.mock.pstmn.io/load-tester/api/v1/start
func Init() {
	log.Println("Webserver init")
	mux := http.NewServeMux()
	mux.HandleFunc("/load-tester/api/v1/start", startLoad)
	handler := cors.AllowAll().Handler(mux)

	if err := http.ListenAndServe(":8080", handler); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
}

func startLoad(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		//data, err := io.ReadAll(req.Body)
		//decoder := json.NewDecoder(data)
		data, byteReadingError := io.ReadAll(req.Body)
		if byteReadingError != nil {
			http.Error(w, fmt.Sprintf("Error while decode requestBody: %s", byteReadingError), http.StatusBadRequest)
			log.Println(fmt.Sprintf("Error while decode requestBody: %s", byteReadingError))
			return
		}
		//decoder := json.NewDecoder(req.Body)
		var body model.StartLoadRequest
		//err := decoder.Decode(&body)
		err := json.Unmarshal(data, &body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while decode requestBody: %s", err), http.StatusBadRequest)
			log.Println(fmt.Sprintf("Error while decode requestBody: %s", err))
			return
		}
		log.Printf("Input request. Url: %s, Body: %s\n", req.URL.Path, body)
		go func() {
			loader.StartLoadTesting(body)
		}()
		response := "Test started"
		getOkResponse(w, response)
	} else {
		http.Error(w, "Expected Post method", http.StatusMethodNotAllowed)
		log.Println("Expected Post method /load-tester/api/v1/start")
		return
	}
	//else if req.Method == http.MethodOptions {
	//	w.Header().Set("Access-Control-Allow-Origin", "127.0.0.1:5173")
	//	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT")
	//	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	//	w.WriteHeader(http.StatusOK)
	//	return
	//}

}

func getOkResponse(w http.ResponseWriter, response string) {
	jsonResp, err := json.Marshal(response)
	if err != nil {

		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
