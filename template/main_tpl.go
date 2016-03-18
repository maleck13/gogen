package template

const(
	FILE_MAIN = "main.go"
	FILE_API_ROUTER = "api/router.go"
	FILE_API_INDEX = "api/indexRoute.go"
	FILE_ERROR_HANDLER = "api/routeErrorHandler.go"
	FILE_ERRORS = "api/errors.go"
	FILE_CMD_SERVER = "cmd/server.go"
	FILE_API_INDEX_TEST = "api/indexRoute_test.go"
	FILE_CONFIG_EXAMPLE = "config/conf.json"
	FILE_CONFIG = "config/config.go"
)

var TEMPLATE_FILES = []string{FILE_MAIN, FILE_API_ROUTER,
	FILE_API_INDEX,FILE_ERROR_HANDLER,FILE_ERRORS,FILE_CMD_SERVER,
FILE_API_INDEX_TEST,FILE_CONFIG,FILE_CONFIG_EXAMPLE}

func GetContent(file string) string {
	if file == FILE_MAIN{
		return getMain()
	}else if file == FILE_API_ROUTER{
		return getApiRouter()
	}else if file == FILE_API_INDEX{
		return getApiIndex()
	}else if file == FILE_CMD_SERVER{
		return getCmdServer()
	}else if file == FILE_ERROR_HANDLER{
		return getApiErrorHandler()
	}else if file == FILE_ERRORS{
		return getApiErrors()
	}else if file == FILE_API_INDEX_TEST{
		return getApiIndexTest()
	}else if file == FILE_CONFIG{
		return getConfig()
	}else if file == FILE_CONFIG_EXAMPLE{
		return getConfigExample()
	}
	return ""

}

func getConfigExample()string{
	return `{"example":"value"}`
}

func getConfig()string{
	return `package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/Sirupsen/logrus"
)


//use an interface to limit access to the config object to read only
type Configuration interface {
	GetExample()string
}

type config struct {
	Example string
}

func (c *config)GetExample()string{
	return c.Example
}



var Conf Configuration

func SetGlobalConfig(path string){
	Conf = &config{}
	file, err := os.Open(path)
	if nil != err{
		logrus.Panic("failed to open config file ", err)
		return;
	}
	defer file.Close()
	data,err := ioutil.ReadAll(file)
	if nil != err{
		logrus.Panic("failed to read config file ", err)
		return;
	}
	if err = json.Unmarshal(data,Conf); err != nil{
		logrus.Panic("failed to decode config file ", err)
		return;
	}
}
`
}


func getApiIndexTest()string{
	return `package api_test

import (
	"testing"
	"net/http/httptest"
	"{basePackage}/api"
	"net/http"
	"github.com/stretchr/testify/assert"
	"{basePackage}/config"
	"io/ioutil"
	"encoding/json"
)

func TestIndexRoute(t *testing.T){
	config.SetGlobalConfig("../config/conf.json")
	server:=httptest.NewServer(api.NewRouter())
	defer server.Close()
	resp, err := http.Get(server.URL)
	assert.NoError(t,err,"did not expect an error")
	assert.Equal(t,200,resp.StatusCode,"expected 200 status code")
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t,err,"did not expect an error reading body")
	data := make(map[string]string)
	err = json.Unmarshal(body,&data)
	assert.NoError(t,err,"did not expect an error reading body")
	if v,ok := data["example"]; ok{
		assert.Equal(t,"value",v,"expected values to match")

	}else{
		assert.Fail(t,"expected returned json to have example key")
	}


}
`
}

func getMain()string{
	return `

package main
import (
	"github.com/codegangsta/cli"
	"os"
	"{basePackage}/cmd"
)

func main()  {

	app := cli.NewApp()
	app.Name = "change_me"
	commands := []cli.Command{
	 cmd.ServeCommand(),
	}
	app.Commands = commands
	app.Run(os.Args)

}

`
}

func getCmdServer()string{
  return	`
  package cmd

import (
	"github.com/codegangsta/cli"
	"{basePackage}/api"
	"github.com/Sirupsen/logrus"
  	"net/http"
	"{basePackage}/config"
)

var port, configPath string
func ServeCommand() cli.Command{
	return cli.Command{
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "port",
				Value:       ":3000",
				Usage:       "serves up the json data",
				Destination: &port,
			},
			cli.StringFlag{
				Name:        "config",
				Value:       "./config/conf.json",
				Usage:       "config file location",
				Destination: &configPath,
			},
		},
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "start the httpgit web service",
		Action:  serve,
	}
}


func serve(context *cli.Context){
	config.SetGlobalConfig(configPath)
	router := api.NewRouter()
	logrus.Info("starting " + context.App.Name + " port "+ port)
	if err := http.ListenAndServe(port, router); err != nil {
		logrus.Fatal(err)
	}
}

`
}


func getApiRouter()string{
	return `package api

import (
	"net/http"
	"github.com/gorilla/mux"
)

type HttpHandler func(wr http.ResponseWriter, req *http.Request) HttpError

func NewRouter()http.Handler{
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/", RouteErrorHandler(IndexHandler)).Methods("GET")
	return router
}

`
}




func getApiIndex()string{
	return `package api

import (
	"net/http"
	"encoding/json"
	"{basePackage}/config"
)

//Example route handler
func IndexHandler(rw http.ResponseWriter, req *http.Request)HttpError{
	encoder := json.NewEncoder(rw);
	data := make(map[string]string)
	data["example"]  = config.Conf.GetExample()
	if err:=encoder.Encode(data); err != nil{
		return NewHttpError(err,http.StatusInternalServerError);
	}
	return nil
}`
}

func getApiErrors()string{
	return `package api

//generic http error wrapper
type  HttpError interface {
	HttpErrorCode() int
}

type HttpHandlerError struct {
	Code int
	Message string
}

func NewHttpError(err error,code int)*HttpHandlerError{
	return &HttpHandlerError{Message:err.Error(),Code:code}
}

func (he *HttpHandlerError)Error()string{
	return he.Message
}

func (he *HttpHandlerError)HttpErrorCode()int{
	return he.Code
}
`
}

func getApiErrorHandler()string{
	return `package api

import (
	"net/http"
	"encoding/json"
)


//Wraps route handlers so that if there is an error returned we don't need to duplicate the error handling
func RouteErrorHandler(handler HttpHandler)http.HandlerFunc {

	return func(wr http.ResponseWriter, req *http.Request) {
		encoder := json.NewEncoder(wr)
		//may change to use a context object containing other data
		if err := handler(wr, req); err != nil {
			wr.WriteHeader(err.HttpErrorCode())
			encoder.Encode(err)
			return;
		}

	}
}



`
}
