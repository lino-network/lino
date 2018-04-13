package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/tx/validator/model"
	"github.com/lino-network/lino/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	log.Println("serve...")
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", "login.html")

	node := rpcclient.NewHTTP("tcp://localhost:46657", "/websocket")

	path := fmt.Sprintf("/%s/key", types.ValidatorKVStoreKey)
	opts := rpcclient.ABCIQueryOptions{
		Height:  0,
		Trusted: true,
	}
	result, err := node.ABCIQueryWithOptions(path, model.GetValidatorListKey(), opts)
	if err != nil {
		log.Println("query failed")
		return
	}
	resp := result.Response
	if resp.Code != uint32(0) {
		log.Println("response not good")
		return
	}
	validatorList := new(model.ValidatorList)
	cdc := app.MakeCodec()
	if err := cdc.UnmarshalJSON(resp.Value, validatorList); err != nil {
		log.Println("unmarshal failed")
		return
	}
	var oncallList = make([]string, len(validatorList.OncallValidators))
	for i, val := range validatorList.OncallValidators {
		oncallList[i] = string(val)
	}
	data := struct {
		Title            string
		OnCallValidators []string
	}{
		Title:            "title",
		OnCallValidators: oncallList,
	}

	varmap := map[string]interface{}{
		"var1":             "value",
		"OnCallValidators": oncallList,
	}
	fmt.Println(data)
	tmpl, _ := template.ParseFiles(lp, fp)
	tmpl.ExecuteTemplate(w, "layout", varmap)
}
