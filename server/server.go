package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	FlagNode    = "node"
	FlagChainID = "chain-id"
)

// SendTxCommand will create a send tx and sign it with the given key
var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "lino server is local server used to interact with blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		fs := http.FileServer(http.Dir("static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))

		http.HandleFunc("/", serveMainPanel)
		http.HandleFunc("/login", serveLogin)

		log.Println("Listening...")
		http.ListenAndServe(":3000", nil)
	},
}

func main() {
	rootCmd.PersistentFlags().StringP(FlagNode, "n", "tcp://localhost:46657", "Node to connect to")
	rootCmd.PersistentFlags().StringP(FlagChainID, "c", "", "ID of chain we connect to")
	viper.BindPFlag(FlagNode, rootCmd.PersistentFlags().Lookup(FlagNode))
	viper.BindPFlag(FlagChainID, rootCmd.PersistentFlags().Lookup(FlagChainID))
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func serveLogin(w http.ResponseWriter, r *http.Request) {
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
	varmap := map[string]interface{}{
		"var1":             "value",
		"OnCallValidators": oncallList,
	}
	tmpl, _ := template.ParseFiles(fp)
	tmpl.ExecuteTemplate(w, "login", varmap)
}

func serveMainPanel(w http.ResponseWriter, r *http.Request) {
	fp := filepath.Join("templates", "index.html")

	log.Println("serve index")
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
	varmap := map[string]interface{}{
		"var1":             "value",
		"OnCallValidators": oncallList,
	}
	tmpl, _ := template.ParseFiles(fp)
	tmpl.ExecuteTemplate(w, "mainDashboard", varmap)
}
