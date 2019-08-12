package restAPI

/**
func GetPagesEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range pages {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Page{})
}

func GetPageEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(pages)
}

func CreatePageEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var page Page
	_ = json.NewDecoder(req.Body).Decode(&page)
	page.ID = params["id"]
	pages = append(pages, page)
	json.NewEncoder(w).Encode(pages)
}

func DeletePageEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range pages {
		if item.ID == params["id"] {
			pages = append(pages[:index], pages[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(pages)
}
**/
