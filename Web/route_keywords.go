package main

import (
	"fmt"
	"html/template"
    "net/http"
    // "io/ioutil"
    "encoding/json"

	"./data"
)

func keywords(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Generating HTML for keywords...")
	sess, err := session(r)
	user, err := data.UserById(sess.UserId)
	if err != nil {
		danger(err, "Cannot find user")
	}

	templates := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/private.navigation.html",
			"templates/keywords.html"))

	type TempStruct struct {
		User    data.User
		Targets []string
	}

	struct_targets, err := user.UsersTargetsByUser()

	var arr_targets []string
	for _, v := range struct_targets {
		arr_targets = append(arr_targets, v.Url)
	}

	infos := TempStruct{user, arr_targets}
    templates.ExecuteTemplate(w, "layout", infos)
    
    _ = err
}

func putKeywordsTargets(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Starting putKeywordsTargets...")
 
    type TempResponse struct {
        SelectedTargets []string `json:"selectedTargets"` // HERE RECEIVE data.Target
        NewKeyword  data.Keyword `json:"newKeyword"`
    }

    response := TempResponse{}

    err := json.NewDecoder(r.Body).Decode(&response)

    fmt.Println(response)

    _ = err

    // Check if keyword exists, if not create, if yes, get keyword id
    // Get all the data of the actual selected targets
    // if the relation userid, targetid, keywordid does not exist, create it, otherwise return a proper message

}