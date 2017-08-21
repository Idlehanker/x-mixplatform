package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"idlehanker.com/x-mixplatform/golan/model"
	. "idlehanker.com/x-mixplatform/golan/util"
)

// var Db *sql.DB

func main() {

	// config := util.Configuration{}

	// util.P(util.Config)
	// util.P("json content:", util.AppConfigContent)
	// db, err := gorm.Open(mysql, args)
	// P("json value ", configContent)
	// defer model.Db.Close()

	cfg := model.ConfigContent["http"].(map[string]interface{})
	P(cfg)

	httpCfg := HTTPConfig{}
	httpCfg.Address = cfg["Address"].(string)
	httpCfg.Static = cfg["Static"].(string)
	httpCfg.ReadTimeout = int64(cfg["ReadTimeout"].(float64))
	httpCfg.WriteTimeout = int64(cfg["WriteTimeout"].(float64))

	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(httpCfg.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.HandleFunc("/", index)
	// mux.HandleFunc("/login", login)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/authenticate", authenticate)

	P(httpCfg)

	server := &http.Server{
		Addr:              httpCfg.Address,
		Handler:           mux,
		ReadHeaderTimeout: time.Duration(httpCfg.ReadTimeout * int64(time.Second)),
		WriteTimeout:      time.Duration(httpCfg.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes:    1 << 20,
	}

	server.ListenAndServe()
}

func index(writer http.ResponseWriter, request *http.Request) {

	P("index page")
}

func login(writer http.ResponseWriter, request *http.Request) {
	t := parseTemplateFiles("login.layout", "public.navbar", "login")
	t.Execute(writer, nil)
	P("login page")
}

func authenticate(writer http.ResponseWriter, request *http.Request) {

	P("authenticate...")
	err := request.ParseForm()
	if err != nil {

		P("error:", err)
	}

	email := request.PostFormValue("email")
	// user, err := model.UserByEmail(email)
	// P("inputed email is '%s'", email)
	fmt.Printf("inputed email is '%s'\n", email)

	user, err := model.UserByEmail(email)
	if err != nil {

		P("error:", err)
	} else {

		password := request.PostFormValue("password")
		fmt.Printf("inputed password is '%s'\n", password)
		if user.Password == Encrypt(password) {

			P("login success!!!")

		} else {
			http.Redirect(writer, request, "/login", 302)
			P("login failure!!!")
		}
		// P("user", user)
	}

}
func parseTemplateFiles(filenames ...string) (t *template.Template) {

	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))

	return
}

// HTTPConfig is config for http server
type HTTPConfig struct {
	Static       string
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
}
