package service

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jaredwarren/app"
	"github.com/jaredwarren/myadmin/db"
)

// Controller implements the home resource.
type Controller struct {
	Mux  *mux.Router
	wapp *app.Service
	db   *db.Database
}

// Register attach mux to service
func Register(wapp *app.Service) *Controller {
	c := &Controller{
		Mux:  wapp.Mux,
		wapp: wapp,
	}

	m := wapp.Mux
	m.HandleFunc("/", c.RunQuery)
	m.HandleFunc("/close", c.Close)
	m.HandleFunc("/run", c.RunQuery).Methods("POST")
	m.HandleFunc("/run", c.RunQuery).Methods("GET")

	m.HandleFunc("/beautify", c.Beautify).Methods("GET", "POST")

	m.HandleFunc("/login", c.Login).Methods("GET")
	m.HandleFunc("/login/{key}", c.LoginHandler).Methods("GET")
	m.HandleFunc("/login", c.LoginHandler).Methods("POST")

	// TODO: download .sql
	// m.HandleFunc("/download", c.LoginHandler).Methods("GET")

	m.HandleFunc("/logout", c.Logout).Methods("GET")

	m.HandleFunc("/{db}", c.SelectDB).Methods("GET")
	m.HandleFunc("/{db}/run", c.RunQuery).Methods("GET")

	m.HandleFunc("/{host}/{db}/run", c.RunQuery).Methods("GET")

	m.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// fmt.Println("~~~~~~~~", r.URL, r.Method, "~~~~~~~~")
			next.ServeHTTP(w, r)
		})
	})

	// TODO:
	//  - show errors!!!!!!
	//  - create edit row/column
	//  - save/load query (maybe make form for params)
	//  -
	//  - ADD COLOR per host
	//  - make host determin db to use i.e. /{host}/{db path}
	// store query size in cookie
	//
	// if column can't be "searched (where)" manually filter (describe orders)

	return c
}

// Close handler.
func (c *Controller) Close(w http.ResponseWriter, r *http.Request) {
	c.db.Close()
	c.wapp.Exit <- nil
}

// Login show current invoice
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	connections, err := db.FetchAll()
	templates := template.Must(template.ParseFiles("templates/login.html", "templates/base.html"))
	templates.ExecuteTemplate(w, "base", &struct {
		Title       string
		Connections map[string]*db.DSN
		Error       error
		Referer     string
	}{
		Title:       "Login",
		Connections: connections,
		Error:       err,
		Referer:     r.Referer(),
	})
}

// LoginHandler show current invoice
func (c *Controller) LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginHandler -> ")
	err := r.ParseForm()
	if err != nil {
		fmt.Println("  - [E]", err)
		loginError(w, r, err)
		return
	}

	var dbc *db.Database

	vars := mux.Vars(r)
	key, ok := vars["key"]
	if ok {
		dbc, err = db.Fetch(key)
		if err != nil {
			fmt.Println("  - [E]", err)
			loginError(w, r, err)
			return
		}
	} else {
		// Create New db
		dbc, err = db.New(r.Form.Get("username"), r.Form.Get("password"), r.Form.Get("host"), r.Form.Get("port"), r.Form.Get("path"))
		if err != nil {
			fmt.Println("  - [E]", err)
			loginError(w, r, err)
			return
		}
		// TODO: save connection (auto?)
	}

	c.db = dbc

	// ref := r.URL.Query().Get("r")
	// fmt.Println(ref)

	fmt.Println("  - done")
	http.Redirect(w, r, "/", 301)
}

func loginError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Println("[E]:", err)
	connections, _ := db.FetchAll()
	templates := template.Must(template.ParseFiles("templates/login.html", "templates/base.html"))
	templates.ExecuteTemplate(w, "base", &struct {
		Title       string
		Connections map[string]*db.DSN
		Error       error
	}{
		Title:       "Login",
		Connections: connections,
		Error:       err,
	})
}

// Logout show current invoice
func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println(" - close db")
	if c.db != nil {
		c.db.Close()
		c.db = nil
	}

	http.Redirect(w, r, "/login", 301)
}

// RunQuery show current invoice
func (c *Controller) RunQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RunQuery:")

	var err error
	if c.db == nil {
		fmt.Println("  - No db")
		http.Redirect(w, r, "/login", 301)
		return
	}

	vars := mux.Vars(r)
	selectedDB := vars["db"]

	limit := 100
	var result *db.Result

	q := r.URL.Query().Get("query")
	if q == "" {
		err = r.ParseForm()
		if err != nil {
			fmt.Println("  - [E] parse form:", err)
			return
		}
		q = r.Form.Get("query")
	}

	// make sure right db is selected
	c.db.Use(selectedDB)

	if q != "" {
		var query *db.Query
		query, err = db.NewQuery(q)
		if err != nil {
			fmt.Println("  - [E] query parse error:", err)
		}

		// Modify Query based on params
		if query != nil {
			// add sort if any
			query.Sort(r.URL.Query().Get("sortname"), r.URL.Query().Get("sortdir"))
			// add search if any
			query.Search(r.URL.Query().Get("search"), r.URL.Query().Get("searchcol"))
			q = query.String()
		}

		// Limit
		ls := r.URL.Query().Get("limit")
		if ls != "" {
			limit, _ = strconv.Atoi(ls)
		}
		fmt.Println("  - Query: ", q)
		// run query
		result, err = c.db.Query(q)
		if err != nil {
			fmt.Println("  - [E] query run error:", err)
		}
		fmt.Println("    - done:", result.Time)
	}

	if result == nil {
		result = &db.Result{
			Rows: nil,
		}
	}

	c.db.Use(selectedDB)

	output := r.URL.Query().Get("output")
	if output == "json" {
		js, err := json.Marshal(result)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	}

	errs := ""
	if err != nil {
		fmt.Println("  - [E]:", err)
		errs = fmt.Sprintf("%s", err)
	}

	// "building" a ui component would be better i.e. like my goext
	funcMap := template.FuncMap{
		"CanDescribe": func(query string) bool {
			return !strings.HasPrefix(strings.ToLower(query), "describe")
		},
		"Beautify": db.Beautify,
		"dump": func(vari interface{}) string {
			return fmt.Sprintf("(%+v)", vari)
		},
		"isNil": func(value interface{}) bool {
			return value == nil
		},
	}
	tpl, err := template.New("base").Funcs(funcMap).ParseFiles("templates/result.html", "templates/overlay.html", "templates/tree.html", "templates/query.html", "templates/home.html", "templates/base.html")
	if err != nil {
		fmt.Printf("  - [E] parse error:%s\n", err.Error())
	}
	templates := template.Must(tpl, err)
	err = templates.ExecuteTemplate(w, "base", &struct {
		Title      string
		Preview    bool
		Query      string
		Result     *db.Result
		DBStruct   *db.Database
		Error      string
		SelectedDB string
		Limit      int
	}{
		Title:      "Home",
		Preview:    true,
		Query:      q,
		Result:     result,
		DBStruct:   c.db,
		Error:      errs,
		SelectedDB: c.db.Using(),
		Limit:      limit,
	})
	if err != nil {
		fmt.Printf("  - [E] execute error:%s\n", err.Error())
	}
}

// SelectDB show current invoice
func (c *Controller) SelectDB(w http.ResponseWriter, r *http.Request) {
	if c.db == nil {
		fmt.Println("SelectDB:  - no db -> login")
		http.Redirect(w, r, "/login", 301)
		return
	}
	vars := mux.Vars(r)
	db := vars["db"]

	// stupid hack for now, should check that db name is "valid",
	if db == "" || db == "favicon.ico" {
		return
	}
	fmt.Println("SelectDB:", db)

	query := fmt.Sprintf("select * from INFORMATION_SCHEMA.tables where TABLE_SCHEMA = '%s' order by TABLE_NAME asc", db)
	http.Redirect(w, r, fmt.Sprintf("/%s/run?query=%s", db, query), 301)
}

// Struct show current invoice
func (c *Controller) Struct(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("templates/struct.html", "templates/base.html"))
	templates.ExecuteTemplate(w, "base", &struct {
		Title       string
		Connections map[string]*db.DSN
		Error       error
		Referer     string
	}{
		Title:   "Login",
		Error:   nil,
		Referer: r.Referer(),
	})
}

// Beautify clean a query
func (c *Controller) Beautify(w http.ResponseWriter, r *http.Request) {
	var err error
	if c.db == nil {
		fmt.Println("  - No db")
		http.Redirect(w, r, "/login", 301)
		return
	}

	q := r.URL.Query().Get("query")
	if q == "" {
		err = r.ParseForm()
		if err != nil {
			fmt.Println("  - [E] parse form:", err)
			return
		}
		q = r.Form.Get("query")
	}

	clean := db.Beautify(q)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(strings.TrimSpace(clean)))
	return
}
