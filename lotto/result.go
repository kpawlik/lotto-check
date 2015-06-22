package lotto

import (
	"appengine"
	"appengine/datastore"

	//	"strings"
	"time"
)

const (
	RESULT_TABLE_NAME = "Result"
)

type Result struct {
	Lotto    []int
	Plus     []int
	Lucky    []int
	IsPlus   bool
	User     string
	Date     time.Time
	WinLotto []int
	WinPlus  []int
}

func NewResult(lotto *LottoResult, lucky *Lucky) *Result {
	var (
		winLotto, winPlus, plus []int
	)
	winLotto = getIntersect(lotto.Lotto, lucky.Lotto)
	if lucky.Plus {
		winPlus = getIntersect(lotto.Plus, lucky.Lotto)
		plus = lotto.Plus
	}
	return &Result{
		Lotto:    lotto.Lotto,
		Plus:     plus,
		Lucky:    lucky.Lotto,
		IsPlus:   lucky.Plus,
		User:     lucky.User,
		Date:     lotto.Date,
		WinLotto: winLotto,
		WinPlus:  winPlus,
	}
}

func (r Result) Save(c appengine.Context) (err error) {
	var (
		key  *datastore.Key
		keys []*datastore.Key
	)
	q := datastore.NewQuery(RESULT_TABLE_NAME).Filter("User=", r.User).Filter("Date=", r.Date).KeysOnly().Limit(1)
	if keys, err = q.GetAll(c, nil); err != nil {
		return
	}
	if len(keys) > 0 {
		key = keys[0]
	} else {
		key = datastore.NewIncompleteKey(c, RESULT_TABLE_NAME, nil)
	}
	_, err = datastore.Put(c, key, &r)
	return
}

func (r Result) lotto() string {
	if len(r.Lotto) > 0 {
		return numbersToLottoStr(r.Lotto)
	}
	return "-"
}
func (r Result) plus() string {
	if len(r.Plus) > 0 {
		return numbersToLottoStr(r.Plus)
	}
	return "-"
}
func (r Result) winLotto() string {
	if len(r.WinLotto) == 0 {
		return "-"
	}
	return numbersToLottoStr(r.WinLotto)
}
func (r Result) winPlus() string {
	if len(r.WinPlus) == 0 {
		return "-"
	}
	return numbersToLottoStr(r.WinPlus)
}
func (r Result) date() string {
	if r.lotto() == "-" {
		return "-"
	}
	return r.Date.Format(DATE_FORMAT)
}

func LastResult(c appengine.Context, user string) (r Result, err error) {
	var (
		res []Result
	)
	q := datastore.NewQuery(RESULT_TABLE_NAME).Filter("User=", user).Order("-Date").Limit(1)
	if _, err = q.GetAll(c, &res); err != nil {
		return
	}
	if len(res) == 1 {
		r = res[0]
	}
	return
}

func getIntersect(arr1 []int, arr2 []int) []int {
	s1 := NewIntSetFrom(arr1)
	s2 := NewIntSetFrom(arr2)
	i := s1.Intersect(s2)
	return i.Items()
}

func Results(c appengine.Context, user string) (res []Result, err error) {
	q := datastore.NewQuery(RESULT_TABLE_NAME).Filter("User=", user).Order("-Date")
	_, err = q.GetAll(c, &res)
	return
}
