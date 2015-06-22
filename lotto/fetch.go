package lotto

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	LOTTO_RESULT_TABLE_NAME = "LottoResult"
)

// Struct to store lotto results
type LottoResult struct {
	Body      []byte
	Lotto     []int
	Plus      []int
	Date      time.Time
	Error     string
	TimeStamp time.Time
}

// constructor
func NewLottoResult(body []byte, err error) *LottoResult {
	result := &LottoResult{Body: body}
	if err != nil {
		result.Error = err.Error()
	}
	result.readLotto()
	result.readPlus()
	result.readDate()
	return result
}

//	Find numbers describe by re in body
func (r *LottoResult) find(re *regexp.Regexp) (numbers []int) {
	lotto := re.Find(r.Body)
	noExp := regexp.MustCompile(`\d+`)
	noBytes := noExp.FindAll(lotto, -1)
	resultsLenght := len(noBytes)
	numbers = make([]int, resultsLenght)
	for i := 0; i < resultsLenght; i++ {
		if no, err := strconv.Atoi(string(noBytes[i])); err == nil {
			numbers[i] = no
		}
	}
	return
}

// read lotto numbers to slot
func (r *LottoResult) readLotto() {
	re, _ := reExp["lotto"]
	if r.Lotto = r.find(re); len(r.Lotto) == 0 {
		r.Error += "No results found for lotto.\n"
	}
}

// read plus numbers to slot
func (r *LottoResult) readPlus() {
	re, _ := reExp["plus"]
	if r.Plus = r.find(re); len(r.Plus) == 0 {
		r.Error += "No results found for plus.\n"
	}
}

// read date to slot
func (r *LottoResult) readDate() {
	var (
		reRes []string
		err   error
	)
	re, _ := reExp["date"]
	dText := re.Find(r.Body)
	dateExp := regexp.MustCompile(`(\d{2})-(\d{2})-(\d{2})`)
	if reRes = dateExp.FindStringSubmatch(string(dText)); len(reRes) == 0 {
		r.Error += fmt.Sprintf("Error reading date (%s).\n", dText)
	}
	day, _ := strconv.Atoi(reRes[1])
	month, _ := strconv.Atoi(reRes[2])
	year, _ := strconv.Atoi(reRes[3])
	dateStr := fmt.Sprintf("%02d-%02d-20%d", day, month, year)
	if r.Date, err = time.Parse(DATE_FORMAT, dateStr); err != nil {
		r.Error += "Error parsing date: " + err.Error()
	}
}

func (r LottoResult) LottoStr() string {
	return numbersToLottoStr(r.Lotto[:])
}

func (r LottoResult) PlusStr() string {
	return numbersToLottoStr(r.Plus[:])
}

func (r *LottoResult) Save(c appengine.Context) (err error) {
	var (
		key  *datastore.Key
		keys []*datastore.Key
		//results []Results
	)
	q := datastore.NewQuery(LOTTO_RESULT_TABLE_NAME).
		KeysOnly().
		Limit(1).
		Filter("Date=", r.Date)
	if keys, err = q.GetAll(c, nil); err != nil {
		r.Error = err.Error()
		return
	}
	if len(keys) == 1 {
		key = keys[0]
	} else {
		key = datastore.NewIncompleteKey(c, LOTTO_RESULT_TABLE_NAME, nil)
	}
	r.TimeStamp = time.Now()
	_, err = datastore.Put(c, key, r)
	return
}
