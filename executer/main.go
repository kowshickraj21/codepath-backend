package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/models"
	"main/scripts"
	"net/http"
)

func main() {
	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		var req models.Req
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		cases := len(req.Testcases)
		var res []models.ResStatus
		var solved int
		if req.Language == "java"{
		    res, solved, err = scripts.JavaExecuter(req, cases)
		}else if req.Language == "cpp" {
		    res, solved, err = scripts.CppExecuter(req, cases)
		}else{
			http.Error(w, "Invalid Language", http.StatusBadRequest)
		}
		if err != nil {
		    fmt.Println("Error",err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp := models.Response{
			Results: res,
			Solved:  solved,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	fmt.Println("Server is running on port 8800")
	log.Fatal(http.ListenAndServe(":8800", nil))
}