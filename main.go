package main

import (
	"log"
	"net/http"
	"encoding/json"
	openrtb2 "github.com/prebid/openrtb/v17/openrtb2"
	adapters "main.go/adapters"
)

func home(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var bidRequest openrtb2.BidRequest
	err := decoder.Decode(&bidRequest)
	if err != nil {
		panic(err)
	}
	huaweiAdsRequest, _ := adapters.MakeRequest(&bidRequest)
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(huaweiAdsRequest)
}

// func makeRequest(openRTBRequest *openrtb2.BidRequest) {
// 	huaweiAdsRequest, _ := adapters.MakeRequest(openRTBRequest)
// 	fmt.Printf("%+v\n", huaweiAdsRequest)
// }

func handleRequests() {
	http.HandleFunc("/",home)
}

func main() {
	handleRequests()
	log.Fatal(http.ListenAndServe(":8081", nil))
}