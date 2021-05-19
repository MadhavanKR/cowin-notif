package cowin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const CALENDAR_BY_DISTRICT_URI = "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict"

var MONTH_MAPPER = map[string]string{
	"January":   "01",
	"February":  "02",
	"March":     "03",
	"April":     "04",
	"May":       "05",
	"June":      "06",
	"July":      "07",
	"August":    "08",
	"September": "09",
	"October":   "10",
	"November":  "11",
	"December":  "12",
}

func GetCalendarForDistrict(districtCode string) (*CalendarByDistrictResponse, error) {
	year, month, day := time.Now().Date()
	today := fmt.Sprintf("%d-%s-%d", day, MONTH_MAPPER[month.String()], year)
	fmt.Println("checking covaixn availability for: ", today)
	queryParams := url.Values{
		"district_id": {districtCode},
		"date":        {today},
	}
	calendarByDistrictUrl := fmt.Sprintf("%s?%s", CALENDAR_BY_DISTRICT_URI, queryParams.Encode())
	calendarByDistrictRequest, _ := http.NewRequest("GET", calendarByDistrictUrl, nil)
	//calendarByDistrictRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93
    //Safari/537.36")
    calendarByDistrictRequest.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")
	httpClient := http.Client{}
	response, httpErr := httpClient.Do(calendarByDistrictRequest)
	defer response.Body.Close()
	if httpErr != nil {
		fmt.Println("error while querying cowin api: ", httpErr)
		return nil, httpErr
	}
	fmt.Println("successfully completed api call, httpStatus: ", response.Status)
	responseBytes, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		fmt.Println("error while reading calendarByDistrict response: ", readErr)
		return nil, readErr
	}
	var calendarByDistrictResponse CalendarByDistrictResponse
	unmarshalErr := json.Unmarshal(responseBytes, &calendarByDistrictResponse)
	if unmarshalErr != nil {
		fmt.Println("error while unmarshalling calendarByDistrict response: ", unmarshalErr)
		return nil, unmarshalErr
	}
	return &calendarByDistrictResponse, nil
}

func GetCovaxinCenters(vaccine string, calendarResponse *CalendarByDistrictResponse) []Center {
	vaxCenters := make([]Center, 0)
	for _, vaxCenter := range calendarResponse.Centers {
		for _, session := range vaxCenter.Sessions {
			if strings.ToLower(session.Vaccine) == strings.ToLower(vaccine) {
				vaxCenters = append(vaxCenters, vaxCenter)
				break
			}
		}
	}
	return vaxCenters
}

func GetAvailableCenters(centers []Center) []Center {
	availableCenters := make([]Center, 0)
	for _, center := range centers {
		for _, session := range center.Sessions {
			if session.AvailableCapacity > 0 {
				availableCenters = append(availableCenters, center)
				break
			}
		}
	}
	return availableCenters
}
