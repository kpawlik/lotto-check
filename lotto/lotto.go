package lotto

import (
	"appengine"
	"appengine/user"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	LOTTO      = `<div class=\"glowna_wyniki_lotto\">[r\n\s.\w\<\>\"\=\/\-\+\!]*<div class=\"[^wynik].*\">`
	LOTTO_PLUS = `<div class=\"glowna_wyniki_lottoplus\">[r\n\s.\w\<\>\"\=\/\-\+\!]*<div class=\"[^wynik].*\">`
	LOTTO_DATE = `<div class=\"start-wyniki_lotto\">[r\n\s\.\w\<\>\"\=\/\-\+\!\:\;?\,]*</div>`
	NUMBERS_RE = `(\d{0,2})\s*`
	URL        = "http://www.lotto.pl/"
	LOTTO_SIZE = 6
	MIN_WIN_NO = 2
)

var (
	tmpls     *template.Template
	reExp                    = map[string]*regexp.Regexp{}
	numbersRe *regexp.Regexp = regexp.MustCompile(NUMBERS_RE)
)

func init() {
	var (
		err error
	)
	funcMap := template.FuncMap{
		"dateStr": DateStr,
		"numbers": NumbersStr,
	}
	tmpls = template.New("tmpls").Funcs(funcMap)
	if tmpls, err = tmpls.ParseFiles(
		"templates/index.html",
		"templates/history.html",
	); err != nil {
		log.Panic("Error during compile templates: ", err)
		//panic(err)
	}
	// Init regular expression map
	strs := map[string]string{"lotto": LOTTO, "plus": LOTTO_PLUS, "date": LOTTO_DATE}
	for k, v := range strs {
		reExp[k] = regexp.MustCompile(v)
	}
	// init handlers
	http.HandleFunc("/", index)
	http.HandleFunc("/check", check)
	http.HandleFunc("/saveLucky", saveLucky)
	http.HandleFunc("/getLucky", getLucky)
	http.HandleFunc("/history", history)
}

// Do nothing
func index(w http.ResponseWriter, r *http.Request) {
	var (
		luckyNbrs, startDate, endDate, isPlus    string
		lottData, lotto, plus, lottoWin, plusWin string
		lastLucky                                Lucky
		lastRes                                  Result
		err                                      error
	)
	//check(w, r)

	c := appengine.NewContext(r)
	u := user.Current(c)
	url, _ := user.LogoutURL(c, "/")

	if lastLucky, err = getLastLucky(c); err == nil {
		luckyNbrs = lastLucky.lotto()
		startDate = lastLucky.startDate()
		endDate = lastLucky.endDate()
		isPlus = lastLucky.plus()
	} else {
		c.Errorf("Get last lucky error: %s", err)
	}

	if lastRes, err = LastResult(c, u.Email); err != nil {
		c.Errorf("Get last results error: %s", err)
	} else {
		lottData = lastRes.date()
		lotto = lastRes.lotto()
		plus = lastRes.plus()
		lottoWin = lastRes.winLotto()
		plusWin = lastRes.winPlus()
	}

	context := struct {
		LogoutURL  string
		User       *user.User
		ResultDate string
		Lotto      string
		Plus       string
		LottoWin   string
		PlusWin    string
		Lucky      string
		Start      string
		End        string
		IsPlus     string
	}{url, u, lottData, lotto, plus, lottoWin, plusWin, luckyNbrs, startDate, endDate, isPlus}
	tmpls.ExecuteTemplate(w, "index", context)
}

// Do nothing
func check(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	pd, err := getPageBody(c, URL)
	result := NewLottoResult(pd, err)
	if err := result.Save(c); err != nil {
		log.Panic(err)
	}
	checkLucky(c, result, w)
}

func saveLucky(w http.ResponseWriter, r *http.Request) {
	context := make(map[string]string)
	c := appengine.NewContext(r)
	lucky := NewLuckyFromRequest(r)
	lucky.Save(c)
	context["status"] = "ok"
	if len(lucky.Error) != 0 {
		//response = `{status: "Error", msg: "` + lucky.Error + `"}`
		context["status"] = "error"
		context["msg"] = lucky.Error
	}
	js, _ := json.Marshal(context)
	fmt.Fprintf(w, "%s", js)
}

func getLucky(w http.ResponseWriter, r *http.Request) {
	context := make(map[string]interface{})
	c := appengine.NewContext(r)
	lucky, _ := getLastLucky(c)
	context["lucky"] = lucky.lotto()
	context["plus"] = lucky.Plus
	context["startDate"] = lucky.startDate()
	context["endDate"] = lucky.endDate()
	js, _ := json.Marshal(context)
	fmt.Fprintf(w, "%s", js)
}

func checkLucky(c appengine.Context, lotto *LottoResult, w http.ResponseWriter) {
	luckys, err := getActiveLucky(c, lotto.Date)
	fmt.Fprintf(w, "Date: %s\n", lotto.Date)
	fmt.Fprintf(w, "Error: %s\n", err)
	for _, lucky := range luckys {
		fmt.Fprintf(w, "  Rec: %q\n", lucky)
		res := NewResult(lotto, lucky)
		err = res.Save(c)
		fmt.Fprintf(w, "RES: %q\n", res)
		fmt.Fprintf(w, "ERR: %q\n", err)
	}
}

func history(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	url, _ := user.LogoutURL(c, "/")
	data, _ := Results(c, u.Email)
	context := struct {
		LogoutURL string
		User      *user.User
		Data      []Result
	}{url, u, data}
	tmpls.ExecuteTemplate(w, "history", context)
}

func DateStr(date time.Time) string {
	return date.Format(DATE_FORMAT)
}
func NumbersStr(numbers []int) string {
	str := make([]string, 0, 6)
	for _, i := range numbers {
		str = append(str, strconv.Itoa(i))
	}
	return strings.Join(str, " ")
}
