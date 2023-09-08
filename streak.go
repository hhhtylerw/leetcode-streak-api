package leetcodestreakapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Actual response from Leetcode GraphQL API.
// I only care about the streak, so I only parse that and errors.
type LeetcodeResponse struct {
	Data struct {
		MatchedUser struct {
			UserCalendar struct {
				// ActiveYears        []int  `json:"activeYears"`
				Streak int `json:"streak"`
				// TotalActiveDays    int    `json:"totalActiveDays"`
				// DccBadges          []any  `json:"dccBadges"`
				// SubmissionCalendar string `json:"submissionCalendar"`
			} `json:"userCalendar"`
		} `json:"matchedUser"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// Sends request for streak. Returns error if present.
func GetStreak(user string) (int, error) {
	var data = strings.NewReader(`{"query":"\n    query userProfileCalendar($username: String!, $year: Int) {\n  matchedUser(username: $username) {\n    userCalendar(year: $year) {\n      activeYears\n      streak\n      totalActiveDays\n      dccBadges {\n        timestamp\n        badge {\n          name\n          icon\n        }\n      }\n      submissionCalendar\n    }\n  }\n}\n    ","variables":{"username":"` + user + `"},"operationName":"userProfileCalendar"}`)
	req, err := http.NewRequest("POST", "https://leetcode.com/graphql/", data)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:109.0) Gecko/20100101 Firefox/117.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Origin", "https://leetcode.com")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://leetcode.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("TE", "trailers")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var leetcodeResponse LeetcodeResponse
	err = json.Unmarshal(bodyText, &leetcodeResponse)
	if err != nil {
		return 0, err
	}

	if len(leetcodeResponse.Errors) > 0 {
		return 0, fmt.Errorf(leetcodeResponse.Errors[0].Message)
	}

	return leetcodeResponse.Data.MatchedUser.UserCalendar.Streak, nil

}
