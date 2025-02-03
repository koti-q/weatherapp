package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const WEATHER_API_URL = "https://api.openweathermap.org/data/2.5/weather"

var apiKey string

func loadEnv() {
	file, err := os.Open("../../.env")
	if err != nil {
		panic("Error opening .env file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "WEATHER_API_KEY=") {
			apiKey = strings.TrimPrefix(line, "WEATHER_API_KEY=")
			break
		}
	}

	if apiKey == "" {
		panic("WEATHER_API_KEY not found in .env file")
	}
	fmt.Println("API key loaded successfully" + apiKey)
}

func init() {
	loadEnv()
}

func getWeater(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		city = "London"
	}

	url := WEATHER_API_URL + "?q=" + city + "&appid=" + apiKey
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write(body)
}

func main() {
	http.HandleFunc("/api/weather", getWeater)
	http.ListenAndServe(":8080", nil)

}
