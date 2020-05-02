package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"regexp"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// Dateiname der yml
const filenameConfig string = "mrs.yml"
const splitter string = "----------------------"

var regex *regexp.Regexp
var config Config

// Config : Struktur der yaml
type Config struct {
	Server struct {
		Port                       string `yaml:"port"`
		FindAllStringSubmatch      bool   `yaml:"findAllStringSubmatch"`
		FindAllStringSubmatchIndex bool   `yaml:"findAllStringSubmatchIndex"`
	} `yaml:"server"`
	Group []struct {
		Name  string `yaml:"name"`
		Regex string `yaml:"regex"`
	} `yaml:"groups"`
}

// Result : Struktur des Ergebnisses
type Result struct {
	Group   []string   `json:"groups"`
	Match   [][]string `json:"matches"`
	Indexes [][]int    `json:"indexes"`
}

func main() {
	log.Println("Starte MultiRegexSuche")
	log.Println(splitter)
	log.Printf("Konfiguration: %v", filenameConfig)
	config.read()
	log.Printf("Port: %v", config.Server.Port)
	log.Printf("Example curl: %v", "curl -XPOST --data @bigtext.txt http://localhost:8080/analysis")
	log.Println(splitter)
	comboRegex := generateRegex(&config)
	regex = compileRegex(comboRegex)
	startWebserver()
}

func (c *Config) read() *Config {
	log.Printf("Lese %v", filenameConfig)
	yamlFile, err := ioutil.ReadFile(filenameConfig)
	if err != nil {
		log.Printf("yaml " + filenameConfig + " nicht gefunden")
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

// Kombiniert alle Regex in einem Regex
func generateRegex(c *Config) string {
	log.Println("Kombiniere Regex ...")
	var comboRegex string
	//comboRegex += "^"
	for index, group := range c.Group {
		if index > 0 {
			comboRegex += "|"
		}
		comboRegex += "("
		comboRegex += "?P<"
		comboRegex += group.Name
		comboRegex += ">"
		comboRegex += group.Regex
		comboRegex += ")"
	}
	//comboRegex += "$"
	log.Println("Regex fertig kombiniert: ", comboRegex)
	return comboRegex
}

func compileRegex(comboRegex string) *regexp.Regexp {
	log.Println("Kompiliere Regex ...")
	r := regexp.MustCompile(comboRegex)
	log.Println("Regex fertig kompiliert")
	return r
}

func startWebserver() {
	log.Println("Starte Webserver ...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/analysis", analysisHandler)
	log.Fatal(http.ListenAndServe(":"+config.Server.Port, router))
}

func analysisHandler(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	start := time.Now()
	var res Result
	res.Group = regex.SubexpNames()
	inputStream := r.Body
	bodyByte, _ := ioutil.ReadAll(inputStream)
	body := string(bodyByte)
	if config.Server.FindAllStringSubmatch {
		res.Match = regex.FindAllStringSubmatch(body, -1)
	}
	if config.Server.FindAllStringSubmatchIndex {
		res.Indexes = regex.FindAllStringSubmatchIndex(body, -1)
	}
	json, _ := json.Marshal(res)
	diff := time.Since(start)
	log.Printf("Analyse fertig in %v", diff)
	w.Write(json)
}

func addCorsHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Allow-Origin", "*")
}
