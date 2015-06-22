package lotto

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/urlfetch"
)

const (
	DATE_FORMAT = `02-01-2006`
)

// read page body from url address
func getPageBody(c appengine.Context, url string) (body []byte, err error) {
	var (
		resp *http.Response
	)
	client := urlfetch.Client(c)
	if resp, err = client.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

// Slice of integer to string array
func numbersToStrArr(numbers []int) []string {
	spliceSize := len(numbers)
	stringArr := make([]string, spliceSize)
	for i := 0; i < spliceSize; i++ {
		stringArr[i] = strconv.FormatInt(int64(numbers[i]), 10)
	}
	return stringArr
}

// Slice of integer to space separated string
func numbersToLottoStr(ints []int) string {
	return strings.Join(numbersToStrArr(ints), " ")
}

// Return date obj from string repr
func parseDate(strDate string) (date time.Time, err error) {
	var (
		ddLen, year, month, day int
		dd                      []string
	)
	dd = strings.Split(strDate, "-")
	if ddLen = len(dd); ddLen < 2 || ddLen > 3 {
		err = errors.New(fmt.Sprintf("Bad data string '%s'", strDate))
		return
	}
	if ddLen == 2 {
		year = time.Now().Year()
	} else {
		y := dd[2]
		if len(y) == 2 {
			y = fmt.Sprintf("20%s", y)
		}
		if year, err = strconv.Atoi(y); err != nil {
			return
		}
	}
	if month, err = strconv.Atoi(dd[1]); err != nil {
		return
	}
	if day, err = strconv.Atoi(dd[0]); err != nil {
		return
	}
	dateS := fmt.Sprintf("%02d-%02d-%d", day, month, year)
	date, err = time.Parse(DATE_FORMAT, dateS)
	return
}

//Parse string to array of ints. String must be space separated ints.
func getNumbers(str string) (nbrs []int, err error) {
	var (
		no, foundLen int
	)
	found := numbersRe.FindAllString(str, -1)
	if found == nil {
		err = errors.New(fmt.Sprintf("Error parsing expresion '%s' to numbers ", str))
		return
	}
	foundLen = len(found)
	nbrs = make([]int, foundLen)
	for i := 0; i < foundLen; i++ {
		if no, err = strconv.Atoi(strings.TrimSpace(found[i])); err != nil {
			return
		}
		nbrs[i] = no
	}
	return
}
