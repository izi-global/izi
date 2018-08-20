package izifix

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/izi-global/izi/cmd/commands"
	"github.com/izi-global/izi/cmd/commands/version"
	iziLogger "github.com/izi-global/izi/logger"
	"github.com/izi-global/izi/logger/colors"
)

var CmdFix = &commands.Command{
	UsageLine: "fix",
	Short:     "Fixes your application by making it compatible with newer versions of IZIGo",
	Long: `As of {{"IZIGo 1.6"|bold}}, there are some backward compatibility issues.

  The command 'fix' will try to solve those issues by upgrading your code base
  to be compatible  with IZIGo version 1.6+.
`,
}

func init() {
	CmdFix.Run = runFix
	CmdFix.PreRun = func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() }
	commands.AvailableCommands = append(commands.AvailableCommands, CmdFix)
}

func runFix(cmd *commands.Command, args []string) int {
	output := cmd.Out()

	iziLogger.Log.Info("Upgrading the application...")

	dir, err := os.Getwd()
	if err != nil {
		iziLogger.Log.Fatalf("Error while getting the current working directory: %s", err)
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".exe") {
			return nil
		}
		err = fixFile(path)
		fmt.Fprintf(output, colors.GreenBold("\tfix\t")+"%s\n", path)
		if err != nil {
			iziLogger.Log.Errorf("Could not fix file: %s", err)
		}
		return err
	})
	iziLogger.Log.Success("Upgrade Done!")
	return 0
}

var rules = []string{
	"izigo.AppName", "izigo.BConfig.AppName",
	"izigo.RunMode", "izigo.BConfig.RunMode",
	"izigo.RecoverPanic", "izigo.BConfig.RecoverPanic",
	"izigo.RouterCaseSensitive", "izigo.BConfig.RouterCaseSensitive",
	"izigo.IZIGoServerName", "izigo.BConfig.ServerName",
	"izigo.EnableGzip", "izigo.BConfig.EnableGzip",
	"izigo.ErrorsShow", "izigo.BConfig.EnableErrorsShow",
	"izigo.CopyRequestBody", "izigo.BConfig.CopyRequestBody",
	"izigo.MaxMemory", "izigo.BConfig.MaxMemory",
	"izigo.Graceful", "izigo.BConfig.Listen.Graceful",
	"izigo.HttpAddr", "izigo.BConfig.Listen.HTTPAddr",
	"izigo.HttpPort", "izigo.BConfig.Listen.HTTPPort",
	"izigo.ListenTCP4", "izigo.BConfig.Listen.ListenTCP4",
	"izigo.EnableHttpListen", "izigo.BConfig.Listen.EnableHTTP",
	"izigo.EnableHttpTLS", "izigo.BConfig.Listen.EnableHTTPS",
	"izigo.HttpsAddr", "izigo.BConfig.Listen.HTTPSAddr",
	"izigo.HttpsPort", "izigo.BConfig.Listen.HTTPSPort",
	"izigo.HttpCertFile", "izigo.BConfig.Listen.HTTPSCertFile",
	"izigo.HttpKeyFile", "izigo.BConfig.Listen.HTTPSKeyFile",
	"izigo.EnableAdmin", "izigo.BConfig.Listen.EnableAdmin",
	"izigo.AdminHttpAddr", "izigo.BConfig.Listen.AdminAddr",
	"izigo.AdminHttpPort", "izigo.BConfig.Listen.AdminPort",
	"izigo.UseFcgi", "izigo.BConfig.Listen.EnableFcgi",
	"izigo.HttpServerTimeOut", "izigo.BConfig.Listen.ServerTimeOut",
	"izigo.AutoRender", "izigo.BConfig.WebConfig.AutoRender",
	"izigo.ViewsPath", "izigo.BConfig.WebConfig.ViewsPath",
	"izigo.StaticDir", "izigo.BConfig.WebConfig.StaticDir",
	"izigo.StaticExtensionsToGzip", "izigo.BConfig.WebConfig.StaticExtensionsToGzip",
	"izigo.DirectoryIndex", "izigo.BConfig.WebConfig.DirectoryIndex",
	"izigo.FlashName", "izigo.BConfig.WebConfig.FlashName",
	"izigo.FlashSeperator", "izigo.BConfig.WebConfig.FlashSeparator",
	"izigo.EnableDocs", "izigo.BConfig.WebConfig.EnableDocs",
	"izigo.XSRFKEY", "izigo.BConfig.WebConfig.XSRFKey",
	"izigo.EnableXSRF", "izigo.BConfig.WebConfig.EnableXSRF",
	"izigo.XSRFExpire", "izigo.BConfig.WebConfig.XSRFExpire",
	"izigo.TemplateLeft", "izigo.BConfig.WebConfig.TemplateLeft",
	"izigo.TemplateRight", "izigo.BConfig.WebConfig.TemplateRight",
	"izigo.SessionOn", "izigo.BConfig.WebConfig.Session.SessionOn",
	"izigo.SessionProvider", "izigo.BConfig.WebConfig.Session.SessionProvider",
	"izigo.SessionName", "izigo.BConfig.WebConfig.Session.SessionName",
	"izigo.SessionGCMaxLifetime", "izigo.BConfig.WebConfig.Session.SessionGCMaxLifetime",
	"izigo.SessionSavePath", "izigo.BConfig.WebConfig.Session.SessionProviderConfig",
	"izigo.SessionCookieLifeTime", "izigo.BConfig.WebConfig.Session.SessionCookieLifeTime",
	"izigo.SessionAutoSetCookie", "izigo.BConfig.WebConfig.Session.SessionAutoSetCookie",
	"izigo.SessionDomain", "izigo.BConfig.WebConfig.Session.SessionDomain",
	"Ctx.Input.CopyBody(", "Ctx.Input.CopyBody(izigo.BConfig.MaxMemory",
	".UrlFor(", ".URLFor(",
	".ServeJson(", ".ServeJSON(",
	".ServeXml(", ".ServeXML(",
	".ServeJsonp(", ".ServeJSONP(",
	".XsrfToken(", ".XSRFToken(",
	".CheckXsrfCookie(", ".CheckXSRFCookie(",
	".XsrfFormHtml(", ".XSRFFormHTML(",
	"izigo.UrlFor(", "izigo.URLFor(",
	"izigo.GlobalDocApi", "izigo.GlobalDocAPI",
	"izigo.Errorhandler", "izigo.ErrorHandler",
	"Output.Jsonp(", "Output.JSONP(",
	"Output.Json(", "Output.JSON(",
	"Output.Xml(", "Output.XML(",
	"Input.Uri()", "Input.URI()",
	"Input.Url()", "Input.URL()",
	"Input.AcceptsHtml()", "Input.AcceptsHTML()",
	"Input.AcceptsXml()", "Input.AcceptsXML()",
	"Input.AcceptsJson()", "Input.AcceptsJSON()",
	"Ctx.XsrfToken()", "Ctx.XSRFToken()",
	"Ctx.CheckXsrfCookie()", "Ctx.CheckXSRFCookie()",
	"session.SessionStore", "session.Store",
	".TplNames", ".TplName",
	"swagger.ApiRef", "swagger.APIRef",
	"swagger.ApiDeclaration", "swagger.APIDeclaration",
	"swagger.Api", "swagger.API",
	"swagger.ApiRef", "swagger.APIRef",
	"swagger.Infomation", "swagger.Information",
	"toolbox.UrlMap", "toolbox.URLMap",
	"logs.LoggerInterface", "logs.Logger",
	"Input.Request", "Input.Context.Request",
	"Input.Params)", "Input.Params())",
	"httplib.IZIGoHttpSettings", "httplib.IZIGoHTTPSettings",
	"httplib.IZIGoHttpRequest", "httplib.IZIGoHTTPRequest",
	".TlsClientConfig", ".TLSClientConfig",
	".JsonBody", ".JSONBody",
	".ToJson", ".ToJSON",
	".ToXml", ".ToXML",
	"izigo.Html2str", "izigo.HTML2str",
	"izigo.AssetsCss", "izigo.AssetsCSS",
	"orm.DR_Sqlite", "orm.DRSqlite",
	"orm.DR_Postgres", "orm.DRPostgres",
	"orm.DR_MySQL", "orm.DRMySQL",
	"orm.DR_Oracle", "orm.DROracle",
	"orm.Col_Add", "orm.ColAdd",
	"orm.Col_Minus", "orm.ColMinus",
	"orm.Col_Multiply", "orm.ColMultiply",
	"orm.Col_Except", "orm.ColExcept",
	"GenerateOperatorSql", "GenerateOperatorSQL",
	"OperatorSql", "OperatorSQL",
	"orm.Debug_Queries", "orm.DebugQueries",
	"orm.COMMA_SPACE", "orm.CommaSpace",
	".SendOut()", ".DoRequest()",
	"validation.ValidationError", "validation.Error",
}

func fixFile(file string) error {
	rp := strings.NewReplacer(rules...)
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	fixed := rp.Replace(string(content))

	// Forword the RequestBody from the replace
	// "Input.Request", "Input.Context.Request",
	fixed = strings.Replace(fixed, "Input.Context.RequestBody", "Input.RequestBody", -1)

	// Regexp replace
	pareg := regexp.MustCompile(`(Input.Params\[")(.*)("])`)
	fixed = pareg.ReplaceAllString(fixed, "Input.Param(\"$2\")")
	pareg = regexp.MustCompile(`(Input.Data\[\")(.*)(\"\])(\s)(=)(\s)(.*)`)
	fixed = pareg.ReplaceAllString(fixed, "Input.SetData(\"$2\", $7)")
	pareg = regexp.MustCompile(`(Input.Data\[\")(.*)(\"\])`)
	fixed = pareg.ReplaceAllString(fixed, "Input.Data(\"$2\")")
	// Fix the cache object Put method
	pareg = regexp.MustCompile(`(\.Put\(\")(.*)(\",)(\s)(.*)(,\s*)([^\*.]*)(\))`)
	if pareg.MatchString(fixed) && strings.HasSuffix(file, ".go") {
		fixed = pareg.ReplaceAllString(fixed, ".Put(\"$2\", $5, $7*time.Second)")
		fset := token.NewFileSet() // positions are relative to fset
		f, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
		if err != nil {
			panic(err)
		}
		// Print the imports from the file's AST.
		hasTimepkg := false
		for _, s := range f.Imports {
			if s.Path.Value == `"time"` {
				hasTimepkg = true
				break
			}
		}
		if !hasTimepkg {
			fixed = strings.Replace(fixed, "import (", "import (\n\t\"time\"", 1)
		}
	}
	// Replace the v.Apis in docs.go
	if strings.Contains(file, "docs.go") {
		fixed = strings.Replace(fixed, "v.Apis", "v.APIs", -1)
	}
	// Replace the config file
	if strings.HasSuffix(file, ".conf") {
		fixed = strings.Replace(fixed, "HttpCertFile", "HTTPSCertFile", -1)
		fixed = strings.Replace(fixed, "HttpKeyFile", "HTTPSKeyFile", -1)
		fixed = strings.Replace(fixed, "EnableHttpListen", "HTTPEnable", -1)
		fixed = strings.Replace(fixed, "EnableHttpTLS", "EnableHTTPS", -1)
		fixed = strings.Replace(fixed, "EnableHttpTLS", "EnableHTTPS", -1)
		fixed = strings.Replace(fixed, "IZIGoServerName", "ServerName", -1)
		fixed = strings.Replace(fixed, "AdminHttpAddr", "AdminAddr", -1)
		fixed = strings.Replace(fixed, "AdminHttpPort", "AdminPort", -1)
		fixed = strings.Replace(fixed, "HttpServerTimeOut", "ServerTimeOut", -1)
	}
	err = os.Truncate(file, 0)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(fixed), 0666)
}
