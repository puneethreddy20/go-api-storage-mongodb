package main

import (
	"errors"
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"github.com/go-api-storage-mongodb/pkg/mongo"
)

var (
	configFilename = flag.String("config", "config.yaml", "The filename of the configuration")
)

const (
	welcomePattern = "/"

	listTagsPattern  = "/listTags/"
	createUpdateTags = "/create-updateTags/"
	deleteTags       = "/deleteTags/"
	createUser       = "/createUser"
	removeUser       = "/removeUser/"
	getMethod        = "GET"
	postMethod       = "POST"
)

type baseConfig struct {
	HttpAddress    string `yaml:"http_address"`
	StorageUrl     string `yaml:"storage_url"`
	DatabaseName   string `yaml:"database_name"`
	CollectionName string `yaml:"collection_name"`
}

type RuntimeState struct {
	Config  baseConfig
	MongoDB mongo.MongoConfig
	DBMutex sync.Mutex
}

//Parsing the config.yml file and storing all the values in state(RuntimeState type).
func parseConfigfile(configFilename string) (RuntimeState, error) {
	var state RuntimeState

	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		err = errors.New("mising config file failure")
		return state, err
	}

	//ioutil.ReadFile returns a byte slice (i.e)(source)
	source, err := ioutil.ReadFile(configFilename)
	if err != nil {
		err = errors.New("cannot read config file")
		return state, err
	}

	//Unmarshall(source []byte,out interface{})decodes the source byte slice/value and puts them in out.
	err = yaml.Unmarshal(source, &state.Config)

	if err != nil {
		err = errors.New("Cannot parse config file")
		log.Printf("Source=%s", source)
		return state, err
	}
	state.MongoDB.Database = state.Config.DatabaseName
	state.MongoDB.Server = state.Config.StorageUrl
	state.MongoDB.Collection = state.Config.CollectionName

	return state, err
}

func main() {

	flag.Parse()

	state, err := parseConfigfile(*configFilename)
	if err != nil {
		log.Println("Error in parseConfigfile function", err)
		return
	}

	http.Handle(welcomePattern, http.HandlerFunc(state.WelcomeHandler))

	http.Handle(createUser, http.HandlerFunc(state.CreateUser))

	http.Handle(removeUser, http.HandlerFunc(state.DeleteUser))

	http.Handle(createUpdateTags, http.HandlerFunc(state.CreateandUpdateTagsforUser))

	http.Handle(deleteTags, http.HandlerFunc(state.DeleteTagsofaUser))

	http.Handle(listTagsPattern, http.HandlerFunc(state.ListallTagsofaUser))

	log.Fatal(http.ListenAndServe(state.Config.HttpAddress, nil))

}
