package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/tx/validator/model"
	"github.com/lino-network/lino/types"
)

const (
	FlagNodeAddr = "addr"
	FlagChainID  = "chain-id"
)

// SendTxCommand will create a send tx and sign it with the given key
func LocalServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "lino server is local server used to interact with blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			fs := http.FileServer(http.Dir("static"))
			http.Handle("/static/", http.StripPrefix("/static/", fs))

			http.HandleFunc("/", serveTemplate)

			log.Println("Listening...")
			http.ListenAndServe(":3000", nil)
		},
	}
	cmd.Flags().String(FlagNodeAddr, "tcp://localhost:46657", "local node address to interact with blockchain")
	cmd.Flags().String(FlagChainID, "lino", "blockchain identity")
	return cmd
}

func main() {
	if err := LocalServerCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	log.Println("serve...")
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", "login.html")

	res, err := QueryLocalStorage(model.GetValidatorListKey(), types.ValidatorKVStoreKey)
	if err != nil {
		log.Println("query failed")
		return
	}
	validatorList := new(model.ValidatorList)
	cdc := app.MakeCodec()
	if err := cdc.UnmarshalJSON(res, validatorList); err != nil {
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
