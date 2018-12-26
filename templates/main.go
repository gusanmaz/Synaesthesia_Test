package templates

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ResultsPageVars struct {
	Lang          string
	PageHeader    string
	TestName      string
	SinglePercent float64
	MultiPercent  float64
	SubjectRatio  float64
}

type TestTextAjax struct {
	ProgressBar string
	Question    string
	PageTitle   string
	TestNo      int
}

var Globals map[string]map[string]interface{}
var Templates *template.Template

func init() {
	Globals = initGlobalVars()
	www_root := os.Getenv("SYN_WWW_ROOT")
	fmt.Println(www_root)
	fmt.Errorf(www_root)
	//dirs  :=[]string {"www/templates", "www/templates/partials"}
	dirs := []string{www_root + "/templates", www_root + "/templates/partials"}
	var err error
	Templates, err = getTemplates(dirs)
	fmt.Println(err)

}

func NewTestTextAjax(lang string, testNo int, askedChar string) *TestTextAjax {
	vars := new(TestTextAjax)
	if lang == "tr" {
		vars.PageTitle = "Sinestezi Testi"
		vars.TestNo = testNo
		vars.ProgressBar = fmt.Sprintf("%%%v tamamlandı", (testNo-1)*10)
		vars.Question = "Yukarıda kaç adet " + askedChar + " var?"
	}
	if lang == "en" {
		vars.PageTitle = "Synaesthesia Test"
		vars.TestNo = testNo
		vars.ProgressBar = fmt.Sprintf("%v%% completed", (testNo-1)*10)
		vars.Question = "How many " + askedChar + "'s are there above?"
	}
	return vars
}

func initGlobalVars() map[string]map[string]interface{} {
	globalVars := make(map[string]map[string]interface{})
	enMap := make(map[string]interface{})
	trMap := make(map[string]interface{})
	globalVars["en"] = enMap
	globalVars["tr"] = trMap

	enMap["Lang"] = "en"
	trMap["Lang"] = "tr"

	enMap["TestTitle"] = "The Synaesthesia Test"
	trMap["TestTitle"] = "Sinestezi Testi"

	trMap["AgeErrorMessage"] = "0 ile 99 arası bir yaş değeri girilmelidir"
	trMap["AgeLabel"] = "Yaş:"
	trMap["AgePlaceholder"] = "26"
	trMap["EmailLabel"] = "Eposta:"
	trMap["EmailPlaceholder"] = "john@doe.net"
	trMap["EmailPlaceholder"] = ""
	trMap["Explanation"] = "Please provide some personal information!"
	trMap["Female"] = "Kadın"
	trMap["GenderLabel"] = "Cinsiyet:"
	trMap["Male"] = "Erkek"
	trMap["NameLabel"] = "İsim & Soyisim:"
	trMap["NamePlaceholder"] = "John Doe"
	trMap["SubmitAnonymousData"] = "Bilgi paylaşmak istemiyorum"
	trMap["SubmitData"] = "Bilgilerimi gönder"

	enMap["AgeErrorMessage"] = "Age value must be between 0 and 99!"
	enMap["AgeLabel"] = "Age:"
	enMap["AgePlaceholder"] = "26"
	enMap["EmailLabel"] = "Email:"
	enMap["EmailPlaceholder"] = "john@doe.net"
	enMap["Explanation"] = "Please provide some personal information!"
	enMap["Female"] = "Female"
	enMap["GenderLabel"] = "Gender:"
	enMap["Male"] = "Male"
	enMap["NameLabel"] = "Name:"
	enMap["NamePlaceholder"] = "John Doe"
	enMap["SubmitAnonymousData"] = "Do not submit personal data"
	enMap["SubmitData"] = "Submit"

	enMap["StartTestText"] = "Start!"
	trMap["StartTestText"] = "Devam et!"

	enMap["ContinueText"] = "Continue!"
	trMap["ContinueText"] = "Devam et!"

	trMap["EmailPlaceholder"] = ""
	enMap["EmailPlaceholder"] = ""

	trMap["AgePlaceholder"] = ""
	enMap["AgePlaceholder"] = ""

	trMap["NamePlaceholder"] = ""
	enMap["NamePlaceholder"] = ""

	trMap["StartOverText"] = "Yeniden Başla"
	enMap["StartOverText"] = "Start Over"

	return globalVars
}

func getTemplates(templateDirs []string) (*template.Template, error) {
	var allFiles []string
	for _, dir := range templateDirs {
		files2, _ := ioutil.ReadDir(dir)
		for _, file := range files2 {
			filename := file.Name()
			if strings.HasSuffix(filename, ".html") {
				filePath := filepath.Join(dir, filename)
				allFiles = append(allFiles, filePath)
			}
		}
	}
	return template.New("").ParseFiles(allFiles...)
}
