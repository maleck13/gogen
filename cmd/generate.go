package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"
	"github.com/maleck13/gogen/template"
	"strings"
	"io/ioutil"
)

var packageRoot string

func GenerateCommand() cli.Command {
	return cli.Command{
		Name:        "generate",
		Action:      generateServer,
		Usage:       "generate --package=github.com/some/thing",
		Description: "generates a go web app under $GOAPTH/src/<package>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "package",
				Value:       "",
				Usage:       "package below $GOPATH/src",
				Destination: &packageRoot,
			},
		},
	}
}

func generateServer(context *cli.Context) {
	if "" == packageRoot {
		logrus.Error("--package required ", GenerateCommand().Usage)
		return
	}

	logrus.Info("generating app at " + getGoPath() + packageRoot)
	var basePath = getGoPath() + packageRoot
	mkDirs(basePath)
	copyFiles(basePath)
}

func getGoPath() string {
	path, present := os.LookupEnv("GOPATH")
	if !present {
		logrus.Error("could not find $GOPATH value")
		return ""
	}
	return path + "/src/"
}

func mkDirs(basePath string) {
	dirs := []string{basePath, basePath + "/cmd", basePath + "/api",basePath+"/config"}
	for _, d := range dirs {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			logrus.Panic(err)
		}
	}

}

func copyFiles(basePath string) {
	for _,t := range template.TEMPLATE_FILES{
		content := template.GetContent(t)
		content = strings.Replace(content,"{basePackage}",packageRoot,-1)
		filePath := basePath + "/" + t
		logrus.Info("creating file ", filePath)
		if err := ioutil.WriteFile(filePath,[]byte(content),os.ModePerm); err != nil{
			logrus.Panic(err)
		}

	}
}
