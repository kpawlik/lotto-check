package lotto

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"net/http"
	//"appengine/user"
	"strings"
	"time"
)

const (
	LUCKY_TABLE_NAME = "Lucky"
)

type Lucky struct {
	Lotto     []int
	Plus      bool
	StartDate time.Time
	EndDate   time.Time
	User      string
	TimeStamp time.Time
	Error     string
}

func NewLucky(lotto []int, plus bool, startDate, endDate time.Time, user, err string) *Lucky {
	return &Lucky{Lotto: lotto,
		Plus:      plus,
		StartDate: startDate,
		EndDate:   endDate,
		User:      user,
		TimeStamp: time.Now(),
		Error:     err}
}

func NewLuckyFromRequest(r *http.Request) *Lucky {
	var (
		lotto              []int
		errStr             []string
		err                error
		startDate, endDate time.Time
	)
	c := appengine.NewContext(r)
	u := user.Current(c)
	if lotto, err = getNumbers(r.FormValue("lotto")); err != nil {
		errStr = append(errStr, err.Error())
	}
	plus := r.FormValue("plus") == "true"
	if startDate, err = parseDate(r.FormValue("startDate")); err != nil {
		errStr = append(errStr, err.Error())
	}
	if endDate, err = parseDate(r.FormValue("endDate")); err != nil {
		errStr = append(errStr, err.Error())
	}
	lucky := NewLucky(lotto, plus, startDate, endDate, u.Email, strings.Join(errStr, "\n"))
	return lucky
}

func (l *Lucky) Save(c appengine.Context) (err error) {
	var (
		key  *datastore.Key
		keys []*datastore.Key
	)
	u := user.Current(c)
	q := datastore.NewQuery(LUCKY_TABLE_NAME).
		KeysOnly().
		Filter("User=", u.Email).
		Limit(1)

	if keys, err = q.GetAll(c, nil); err != nil {
		l.Error = err.Error()
	}
	if len(keys) == 1 {
		key = keys[0]
	} else {
		key = datastore.NewIncompleteKey(c, LUCKY_TABLE_NAME, nil)
	}
	l.TimeStamp = time.Now()
	_, err = datastore.Put(c, key, l)
	return
}

func (l Lucky) startDate() string {
	if l.lotto() == "-" {
		return "-"
	}
	return l.StartDate.Format(DATE_FORMAT)
}

func (l Lucky) endDate() string {
	if l.lotto() == "-" {
		return "-"
	}
	return l.EndDate.Format(DATE_FORMAT)
}
func (l Lucky) lotto() string {
	if len(l.Lotto) > 0 {
		return numbersToLottoStr(l.Lotto)
	}
	return "-"
}
func (l Lucky) plus() string {
	if l.lotto() == "-" {
		return "-"
	}
	if l.Plus {
		return "Yes"
	} else {
		return "No"
	}
}
func getLastLucky(c appengine.Context) (lucky Lucky, err error) {
	var (
		luckys []Lucky
	)
	u := user.Current(c)
	email := u.Email
	q := datastore.NewQuery(LUCKY_TABLE_NAME).Filter("User=", email).Limit(1)

	if _, err = q.GetAll(c, &luckys); err != nil || len(luckys) == 0 {
		return
	}
	lucky = luckys[0]
	return
}

func getActiveLucky(c appengine.Context, date time.Time) (lucky []*Lucky, err error) {
	var (
		luckys []Lucky
	)
	q := datastore.NewQuery(LUCKY_TABLE_NAME).Filter("EndDate>=", date)
	_, err = q.GetAll(c, &luckys)
	for _, rec := range luckys {
		if rec.StartDate.Equal(date) || rec.StartDate.Before(date) {
			lucky = append(lucky, &rec)
		}
	}
	return
}
