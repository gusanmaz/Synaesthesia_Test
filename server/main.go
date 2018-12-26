package server

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gusanmaz/Synaesthesia_Test/datastore"
	"github.com/gusanmaz/Synaesthesia_Test/rectangle"
	"github.com/gusanmaz/Synaesthesia_Test/stats"
	"github.com/gusanmaz/Synaesthesia_Test/templates"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	Email          string
	Name           string
	Age            int
	Gender         string
	ResponseTime   []int
	Locale         string
	UserAgent      string
	Sid            string
	UnixTime       int64
	TestsGenerated int
}

type ETagData struct {
	ETag string
	When int64
}

type myFileServer struct {
	BasePath  string
	subPath   string
	ETagCache map[string]ETagData
}

type BlogServer struct {
	BasePath  string
	subPath   string
}

type UserResponse struct {
	Correct     bool `json:"correct"`
	Answer      int  `json:"answer""`
	ElapsedTime int  `json:"elapsedTime""`
}

type ResultsPageVars struct {
	Lang          string
	SinglePercent float64
	MultiPercent  float64
	SubjectRatio  float64
}

type AboutTestPageVars struct {
	Lang string
}

var (
	store       *sessions.CookieStore
	fileServer         myFileServer
	blogServer         BlogServer
	sid2Test           map[string]int
	mongo              datastore.Mongo
	staticResourceRoot string
	testPath           string
	blogPath           string
	sessionKey         []byte  // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
)

var (
	testClasses = [][]string{
		{"green2", "gray5"},
		{"gray2", "gray5"},
		{"blue5", "gray2"},
		{"gray5", "gray2"},
		{"green6", "gray8"},
		{"gray6", "gray8"},
		{"blue8", "gray6"},
		{"gray8", "gray6"},
		{"blue4", "greeny"},
		{"gray4", "grayy"},
	}
)

func init() {
	staticResourceRoot = os.Getenv("SYN_WWW_ROOT")
	sessionKey         = []byte(os.Getenv("SYN_SESS_KEY"))
	blogPath           = os.Getenv("BLOG_WWW_ROOT")
	fmt.Println("BlogPath:", blogPath)
	store       = sessions.NewCookieStore(sessionKey)

	testPath = "/st/"
	sid2Test = make(map[string]int, 10)
	fileServer = myFileServer{staticResourceRoot, testPath, map[string]ETagData{}}
	blogServer = BlogServer{blogPath, "/blog/"}
}

func Router(m datastore.Mongo) *mux.Router {
	mongo = m
	mainRouter := mux.NewRouter()

	blogRouter := mainRouter.PathPrefix("/blog/").Subrouter()
	testRouter := mainRouter.PathPrefix(testPath).Subrouter()

	// TODO May be needed in future

	testRouter.HandleFunc("/js/data.json", dataH)
	testRouter.HandleFunc("/js/data_text.json", dataTextH)

	testRouter.HandleFunc("/processUserResponse", processUserResponseH)
	testRouter.HandleFunc("/processLastResponse", processLastResponseH)

	testRouter.HandleFunc("/intro.html", introPageH)
	testRouter.HandleFunc("/index.html", introPageH)
	testRouter.HandleFunc("/", introPageH)

	testRouter.HandleFunc("/results.html", processResultsH)
	testRouter.HandleFunc("/user.html", userPageH)
	testRouter.HandleFunc("/initTest", initTestH)
	testRouter.HandleFunc("/about_test.html", aboutTestH)
	testRouter.HandleFunc("/error.html", errorH)

	testRouter.HandleFunc("/about_test.html", aboutTestH)
	testRouter.HandleFunc("/error.html", errorH)
	testRouter.PathPrefix("/").Handler(fileServer)
	blogRouter.PathPrefix("/").Handler(blogServer)
	return mainRouter
}

func blogH(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This is a blog page!")
}

func indexH(w http.ResponseWriter, r *http.Request) {
	t := templates.Templates
	globals := templates.Globals
	lang := DetectLanguage(r.Header.Get("Accept-Language"))
	globals["en"]["PageTitle"] = "Error"
	globals["tr"]["PageTitle"] = "Hata"
	//TODO Agressive take on reaching an error page. May need to reconsider removal of cookies
	ClearAllCookies(w)
	t.ExecuteTemplate(w, "ErrorHTML", globals[lang])
}

func errorH(w http.ResponseWriter, r *http.Request) {
	t := templates.Templates
	globals := templates.Globals
	lang := DetectLanguage(r.Header.Get("Accept-Language"))
	globals["en"]["PageTitle"] = "Error"
	globals["tr"]["PageTitle"] = "Hata"
	//TODO Agressive take on reaching an error page. May need to reconsider removal of cookies
	ClearAllCookies(w)
	t.ExecuteTemplate(w, "ErrorHTML", globals[lang])
}

func aboutTestH(w http.ResponseWriter, r *http.Request) {
	testNo := 1
	testsGenerated := 1
	maxAge := 10000
	lang := DetectLanguage(r.FormValue("lang"))

	sessionID := GenerateSessionID()
	sidCookie := http.Cookie{Name: "sid", Value: sessionID, Path: "/", HttpOnly: true, MaxAge: int(maxAge)}
	testNoCookie := http.Cookie{Name: "test_no", Value: strconv.FormatInt(int64(testNo), 10), Path: "/", MaxAge: int(maxAge)}
	langCookie := http.Cookie{Name: "lang", Value: lang, Path: "/", MaxAge: int(maxAge)}
	http.SetCookie(w, &sidCookie)
	http.SetCookie(w, &testNoCookie)
	http.SetCookie(w, &langCookie)

	session, err := store.Get(r, "store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	session.Values["lang"] = lang
	session.Values["test_no"] = testNo
	session.Values["tests_generated"] = testsGenerated
	session.Values["elapsed_time"] = "0"
	session.Values["sid"] = sessionID


	sid2Test[sessionID] = testNo

	t := templates.Templates
	globals := templates.Globals
	globals["en"]["PageTitle"] = "About Test Content"
	globals["tr"]["PageTitle"] = "Test İçeriği Hakkında"
	err = session.Save(r, w)
	if err != nil{
		fmt.Println(err)
	}

	t.ExecuteTemplate(w, "AboutTestHTML", globals[lang])
}

func introPageH(w http.ResponseWriter, r *http.Request) {
	lang := DetectLanguage(r.Header.Get("Accept-Language"))
	queryValues := r.URL.Query()
	queryLang := queryValues.Get("lang")
	if queryLang != "" {
		lang = DetectLanguage(queryLang)
	}

	g := templates.Globals
	g["en"]["PageTitle"] = "Introducing Test"
	g["tr"]["PageTitle"] = "Test Hakkında"

	t := templates.Templates
	t.ExecuteTemplate(w, "IntroHTML", g[lang])

}

func userPageH(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	lang := session.Values["lang"].(string)

	g := templates.Globals
	g["en"]["PageTitle"] = "One more thing :)"
	g["tr"]["PageTitle"] = "Son bir soru :)"

	t := *templates.Templates
	t.ExecuteTemplate(w, "UserHTML", templates.Globals[lang])
}

func initTestH(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "store")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	lang := session.Values["lang"].(string)

	t := templates.Templates
	g := templates.Globals

	g["en"]["PageTitle"] = "Progress"
	g["tr"]["PageTitle"] = "Durum"

	t.ExecuteTemplate(w, "TestHTML", g[lang])
}

func dataH(w http.ResponseWriter, r *http.Request) {
	if !IsCookieIntact(r) {
		http.Redirect(w, r, "error.html", 302)
	}

	session, _ := store.Get(r, "store")
	sid := session.Values["sid"].(string)
	testNo := session.Values["test_no"]
	testNoNum := testNo.(int)

	// Probably a redundant test
	internalTestNo, found := sid2Test[sid]
	if !found || internalTestNo != testNo {
		http.Redirect(w, r, "error.html", 302)
	}

	jsonObj := rectangle.GetTestJSON(testClasses[testNoNum-1][0]+".svg", testClasses[testNoNum-1][1]+".svg")
	fmt.Fprint(w, jsonObj)
}

func dataTextH(w http.ResponseWriter, r *http.Request) {
	if !IsCookieIntact(r) {
		http.Redirect(w, r, "error.html", 302)
	}

	session, _ := store.Get(r, "store")
	sid := session.Values["sid"].(string)
	testNo := session.Values["test_no"]
	lang := session.Values["lang"].(string)
	testNoNum := testNo.(int)

	// Probably a redundant test
	internalTestNo, found := sid2Test[sid]
	if !found || internalTestNo != testNo {
		http.Redirect(w, r, "error.html", 302)
	}

	askedClass := testClasses[testNoNum-1][0]
	askedChar := askedClass[len(askedClass)-1:]
	testText := templates.NewTestTextAjax(lang, testNoNum, askedChar)

	jsObj, _ := json.Marshal(&testText)
	fmt.Fprint(w, string(jsObj))
}

func processLastResponseH(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "store")
	if !IsCookieIntact(r) {
		http.Redirect(w, r, "error.html", 302)
	}

	sid := session.Values["sid"].(string)
	testNo := session.Values["test_no"]
	testNoNum := testNo.(int)

	//True normally
	if testNoNum > 10 && sid2Test[sid] > 10 {
		session.Values["test_completed"] = "true"
		session.Save(r, w)
		http.Redirect(w, r, "user.html", 302)
	} else { //Something unexpected happens
		http.Redirect(w, r, "error.html", 302)
	}

}

func processResultsH(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "store")
	testComplete := session.Values["test_completed"]
	if !IsCookieIntact(r) || testComplete == nil || testComplete.(string) != "true" {
		http.Redirect(w, r, "error.html", 302)
	}

	undefinedVal := "na"
	email := undefinedVal
	name := undefinedVal
	gender := undefinedVal
	age := -1
	sid := session.Values["sid"].(string)
	lang := session.Values["lang"].(string)
	testsGenerated := session.Values["tests_generated"].(int)

	if r.Method == "POST" {
		email = r.FormValue("email")
		name = r.FormValue("name")
		age, _ = strconv.Atoi(r.FormValue("age"))
		gender = r.FormValue("gender")

		emailErr := checkmail.ValidateFormat(email)
		ageValid := ValidateAge(age)

		if emailErr != nil || !ageValid {
			//TODO Redirection to more specific error page would be nice
			http.Redirect(w, r, "error.html", 302)
		}
	}

	userAgent := r.Header.Get("User-Agent")
	unixTime := time.Now().Unix()

	elapsedTime := session.Values["elapsed_time"].(string)
	elapsedTime = strings.TrimPrefix(elapsedTime, "0,")

	timeVals := strings.Split(elapsedTime, ",")
	subjectTimes := make([]int, len(timeVals))
	for i, v := range timeVals {
		subjectTimes[i], _ = strconv.Atoi(v)
	}

	result := Result{
		Name:           name,
		Email:          email,
		Age:            age,
		Gender:         gender,
		ResponseTime:   subjectTimes,
		Locale:         lang,
		UserAgent:      userAgent,
		Sid:            sid,
		UnixTime:       unixTime,
		TestsGenerated: testsGenerated,
	}

	var err error
	c := mongo.Session.DB("synaesthesia").C("results")
	err = c.Insert(&result)
	if err != nil {
		log.Fatal(err)
	}

	var population []Result
	query := bson.M{"testsgenerated": bson.M{"$lt": 20}}
	err = c.Find(query).All(&population)
	populationTimes := GetElapsedTimeArr(population)

	s := stats.CalculateUserStatistics(subjectTimes, populationTimes)
	stats.InitStats(s)

	t := templates.Templates
	g := templates.Globals
	t.ExecuteTemplate(w, "ResultsHTML", g[lang])
}

func processUserResponseH(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var userResponse UserResponse
	err := decoder.Decode(&userResponse)
	if err != nil {
		panic(err)
	}

	correct := userResponse.Correct
	elapsedTime := userResponse.ElapsedTime

	session, _ := store.Get(r, "store")
	defer session.Save(r, w)
	if !IsCookieIntact(r) {
		http.Redirect(w, r, "error.html", 302)
	}

	testNo := session.Values["test_no"]
	testNoNum := testNo.(int)
	testsGenerated := session.Values["tests_generated"]
	testsGeneratedNum := testsGenerated.(int)
	maxAge := 10000

	session.Values["tests_generated"] = testsGeneratedNum + 1
	if correct {
		session.Values["test_no"] = testNoNum + 1
		sid := session.Values["sid"].(string)
		sid2Test[sid] = testNoNum + 1
		timeArr := session.Values["elapsed_time"].(string)
		session.Values["elapsed_time"] = fmt.Sprintf("%v,%v", timeArr, elapsedTime)
		cookie := http.Cookie{Name: "test_no", Value: strconv.FormatInt(int64(testNoNum+1), 10), Path: "/", MaxAge: int(maxAge)}
		http.SetCookie(w, &cookie)

		// FIXME This test has no use
		if testNoNum > 10 {
			http.Redirect(w, r, "error.html", 302)
		}
	}
}

func (fs myFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, fs.subPath)
	filePath := fs.BasePath + "/" + path
	fmt.Println(filePath)
	fileStats, err := os.Stat(filePath)

	if strings.Contains(r.URL.Path, "/templates/") {
		FileNotFound404H(w, r, http.StatusNotFound)
		return
	} else if os.IsNotExist(err) {
		//TODO A more specific permission denied page may be considered
		FileNotFound404H(w, r, http.StatusNotFound)
		return
	} else {
		// TODO: Informing through a log mechanism would be nicer
		log.Println("Serving file at: " + filePath)

		etagData, exists := fs.ETagCache[filePath]
		etagVal := etagData.ETag
		fileModTime := fileStats.ModTime().Unix()
		if !exists || fileModTime != etagData.When {
			dat, err := ioutil.ReadFile(filePath)
			if err != nil {
				panic(err)
			}
			log.Println("New map entry created")

			etagVal = fmt.Sprintf("%x", md5.Sum([]byte(dat)))
			fs.ETagCache[filePath] = ETagData{etagVal, fileModTime}
		}

		log.Println("ETag value is set to: ", etagVal)
		w.Header().Set("ETag", etagVal)

		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, etagVal) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		} else {
			var maxAge int64 = 30 * 60 * 60
			w.Header().Set("Cache-Control", "max-age="+strconv.FormatInt(maxAge, 10)+", public, must-revalidate")
			w.Header().Set("ETag", etagVal)
			http.ServeFile(w, r, filePath)
		}
	}
}

func (fs BlogServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I am at blog")
	fmt.Println(r.URL.Path)
	path := strings.TrimPrefix(r.URL.Path, fs.subPath)
	filePath := fs.BasePath + "/" + path
	fmt.Println(filePath)
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		//TODO A more specific permission denied page may be considered
		FileNotFound404H(w, r, http.StatusNotFound)
		return
	} else {
			// TODO: Informing through a log mechanism would be nicer
			log.Println("Serving file at: " + filePath)
			http.ServeFile(w, r, filePath)
	}
}

func FileNotFound404H(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		session, _ := store.Get(r, "store")
		langSess := session.Values["lang"]
		lang := DetectLanguage(r.Header.Get("User-Agent"))
		if langSess != nil {
			lang = langSess.(string)
		}

		t := templates.Templates
		globals := templates.Globals
		globals["en"]["PageTitle"] = "404 Error"
		globals["tr"]["PageTitle"] = "404 Hatası"
		t.ExecuteTemplate(w, "File404HTML", globals[lang])
	}
}

func ClearAllCookies(w http.ResponseWriter) {
	DeleteCookieHandler(w, "sid")
	DeleteCookieHandler(w, "store")
	DeleteCookieHandler(w, "test_no")
	DeleteCookieHandler(w, "test_complete")
}

func IsCookieIntact(r *http.Request) bool {
	session, _ := store.Get(r, "store")
	sid, _ := r.Cookie("sid")
	return sid.Value == session.Values["sid"]
}

func GetElapsedTimeArr(res []Result) [][]int {
	sliceLen := len(res)
	timeVals := make([][]int, sliceLen)
	for i, _ := range res {
		timeVals[i] = res[i].ResponseTime
	}
	return timeVals
}
