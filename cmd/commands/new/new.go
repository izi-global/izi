// Copyright 2018 IZI Global
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package new

import (
	"fmt"
	"os"
	path "path/filepath"
	"strings"

	"github.com/izi-global/izi/cmd/commands"
	"github.com/izi-global/izi/cmd/commands/version"
	iziLogger "github.com/izi-global/izi/logger"
	"github.com/izi-global/izi/logger/colors"
	"github.com/izi-global/izi/utils"
)

var CmdNew = &commands.Command{
	UsageLine: "new [appname]",
	Short:     "Creates a IZIGo application",
	Long: `
Creates a IZIGo application for the given app name in the current directory.

  The command 'new' creates a folder named [appname] and generates the following structure:

            ├── main.go
            ├── {{"conf"|foldername}}
            │     └── app.conf
            ├── {{"controllers"|foldername}}
            │     └── default.go
            ├── {{"models"|foldername}}
            ├── {{"routers"|foldername}}
            │     └── router.go
            ├── {{"tests"|foldername}}
            │     └── default_test.go
            ├── {{"static"|foldername}}
            │     └── {{"js"|foldername}}
            │     └── {{"css"|foldername}}
            │     └── {{"img"|foldername}}
            └── {{"views"|foldername}}
                  └── index.tpl

`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    CreateApp,
}

var appconf = `appname = {{.Appname}}
httpport = 8080
runmode = dev
`

var maingo = `package main

import (
	_ "{{.Appname}}/routers"
	"github.com/izi-global/izigo"
)

func main() {
	izigo.Run()
}

`
var router = `package routers

import (
	"{{.Appname}}/controllers"
	"github.com/izi-global/izigo"
)

func init() {
    izigo.Router("/", &controllers.MainController{})
}
`

var test = `package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"runtime"
	"path/filepath"
	_ "{{.Appname}}/routers"

	"github.com/izi-global/izigo"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	izigo.TestIZIGoInit(apppath)
}


// TestIZIGo is a sample to run an endpoint test
func TestIZIGo(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	izigo.IZIApp.Handlers.ServeHTTP(w, r)

	izigo.Trace("testing", "TestIZIGo", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	                So(w.Code, ShouldEqual, 200)
	        })
	        Convey("The Result Should Not Be Empty", func() {
	                So(w.Body.Len(), ShouldBeGreaterThan, 0)
	        })
	})
}

`

var controllers = `package controllers

import (
	"github.com/izi-global/izigo"
)

type MainController struct {
	izigo.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "go.izi.asia"
	c.Data["Email"] = "dotiendiep@gmail.com"
	c.TplName = "index.tpl"
}
`

var indextpl = `<!DOCTYPE html>

<html>
<head>
  <title>IZIGo</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <link rel="shortcut icon" href="data:image/png;base64,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" type="image/x-icon" />

  <style type="text/css">
    *,body {
      margin: 0px;
      padding: 0px;
    }

    body {
      margin: 0px;
      font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
      font-size: 14px;
      line-height: 20px;
      background-color: #fff;
    }

    header,
    footer {
      width: 960px;
      margin-left: auto;
      margin-right: auto;
    }

    .logo {
      background-image: url('data:image/png;base64,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');
      background-repeat: no-repeat;
      -webkit-background-size: 100px 100px;
      background-size: 100px 100px;
      background-position: center center;
      text-align: center;
      font-size: 42px;
      padding: 250px 0 70px;
      font-weight: normal;
      text-shadow: 0px 1px 2px #ddd;
    }

    header {
      padding: 100px 0;
    }

    footer {
      line-height: 1.8;
      text-align: center;
      padding: 50px 0;
      color: #999;
    }

    .description {
      text-align: center;
      font-size: 16px;
    }

    a {
      color: #444;
      text-decoration: none;
    }

    .backdrop {
      position: absolute;
      width: 100%;
      height: 100%;
      box-shadow: inset 0px 0px 100px #ddd;
      z-index: -1;
      top: 0px;
      left: 0px;
    }
  </style>
</head>

<body>
  <header>
    <h1 class="logo">Welcome to IZIGo</h1>
    <div class="description">
      IZIGo is a simple & powerful Go web framework which is inspired by tornado and sinatra.
    </div>
  </header>
  <footer>
    <div class="author">
      Official website:
      <a href="http://{{.Website}}">{{.Website}}</a> /
      Contact me:
      <a class="email" href="mailto:{{.Email}}">{{.Email}}</a>
    </div>
  </footer>
  <div class="backdrop"></div>

  <script src="/static/js/reload.min.js"></script>
</body>
</html>
`

var reloadJsClient = `function b(a){var c=new WebSocket(a);c.onclose=function(){setTimeout(function(){b(a)},2E3)};c.onmessage=function(){location.reload()}}try{if(window.WebSocket)try{b("ws://localhost:12450/reload")}catch(a){console.error(a)}else console.log("Your browser does not support WebSockets.")}catch(a){console.error("Exception during connecting to Reload:",a)};
`

func init() {
	commands.AvailableCommands = append(commands.AvailableCommands, CmdNew)
}

func CreateApp(cmd *commands.Command, args []string) int {
	output := cmd.Out()
	if len(args) != 1 {
		iziLogger.Log.Fatal("Argument [appname] is missing")
	}

	apppath, packpath, err := utils.CheckEnv(args[0])
	if err != nil {
		iziLogger.Log.Fatalf("%s", err)
	}

	if utils.IsExist(apppath) {
		iziLogger.Log.Errorf(colors.Bold("Application '%s' already exists"), apppath)
		iziLogger.Log.Warn(colors.Bold("Do you want to overwrite it? [Yes|No] "))
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}

	iziLogger.Log.Info("Creating application...")

	os.MkdirAll(apppath, 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "conf"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "conf")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "controllers"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "controllers")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "models"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "models")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "routers"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "routers")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "tests"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "tests")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "static"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "static")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "static", "js"), 0755)
	utils.WriteToFile(path.Join(apppath, "static", "js", "reload.min.js"), reloadJsClient)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "static", "js")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "static", "css"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "static", "css")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "static", "img"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "static", "img")+string(path.Separator), "\x1b[0m")
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "views")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "views"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "conf", "app.conf"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "conf", "app.conf"), strings.Replace(appconf, "{{.Appname}}", path.Base(args[0]), -1))

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "controllers", "default.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "controllers", "default.go"), controllers)

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "views", "index.tpl"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "views", "index.tpl"), indextpl)

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "routers", "router.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "routers", "router.go"), strings.Replace(router, "{{.Appname}}", packpath, -1))

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "tests", "default_test.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "tests", "default_test.go"), strings.Replace(test, "{{.Appname}}", packpath, -1))

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "main.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "main.go"), strings.Replace(maingo, "{{.Appname}}", packpath, -1))

	iziLogger.Log.Success("New application successfully created!")
	return 0
}
