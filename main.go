package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Location struct {
	Name    string  `json:"name"`
	Region  string  `json:"region"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type Current struct {
	Temp_c     float64    `json:"temp_c"`
	Temp_f     float64    `json:"temp_f"`
	AirQuality AirQuality `json:"air_quality"`
	Condition  Condition  `json:"condition"`
}

type AirQuality struct {
	Co           float64 `json:"co"`
	No2          float64 `json:"no2"`
	O3           float64 `json:"o3"`
	So2          float64 `json:"so2"`
	Pm2_5        float64 `json:"pm2_5"`
	Pm10         float64 `json:"pm10"`
	UsEpAqi      float64 `json:"us-epa-index"`
	GbDefraIndex float64 `json:"gb-defra-index"`
}

type Condition struct {
	Text string  `json:"text"`
	Icon string  `json:"icon"`
	Code float64 `json:"code"`
}
type DailyHour struct {
	Time       string     `json:"time"`
	TImeEpoch  int        `json:"time_epoch"`
	Condition  Condition  `json:"condition"`
	Temp_c     float64    `json:"temp_c"`
	WillItRain float64    `json:"will_it_rain"`
	AirQuality AirQuality `json:"air_quality"`
}

type Forecastday struct {
	Date       string      `json:"date"`
	Date_epoch int         `json:"date_epoch"`
	Hour       []DailyHour `json:"hour"`
}

type Forecast struct {
	Forecastday []Forecastday `json:"forecastday"`
}

type Weather struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
	Forecast Forecast `json:"forecast"`
}

func main() {
	res, err := http.Post("http://api.weatherapi.com/v1/forecast.json?key=231b5de02cd241bbb2390630250701&q=Pune&days=1&aqi=yes", "application/json", nil)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error ")
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Println("Error in unmarshalling")
	}
	location, current, hours, date := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour, weather.Forecast.Forecastday[0].Date

	fmt.Printf("%s, %s, %.0f, %s \n", location.Name, location.Country, current.Temp_c, current.Condition.Text)
	fmt.Printf("Date: %v \n", date)
	fmt.Println("Time, Temp, Condition, Will it rain, PM2.5, PM10, US EPA AQI")
	for _, hour := range hours {
		date := time.Unix(int64(hour.TImeEpoch), 0)
		fmt.Printf("%s, %.0f, %s, %.0f, %.0f, %.0f, %.0f \n", date.Format("15:04"), hour.Temp_c, hour.Condition.Text, hour.WillItRain, hour.AirQuality.Pm2_5, hour.AirQuality.Pm10, hour.AirQuality.UsEpAqi)
	}
}
