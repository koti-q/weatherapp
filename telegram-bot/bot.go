package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var token string
var chatIDs = make(map[int64]bool)

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func loadEnv() {
	file, err := os.Open("../config/.env")
	if err != nil {
		panic("Error opening .env file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "TG_BOT_TOKEN=") {
			token = strings.TrimPrefix(line, "TG_BOT_TOKEN=")
			break
		}
	}

	if token == "" {
		panic("TG_BOT_TOKEN not found in .env file")
	}
	log.Println("Token loaded successfully" + token)
}

func getWeatherUpdate(city string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/weather?city=%s", city))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var weather WeatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return "", err
	}

	tempCelsius := weather.Main.Temp - 273.15
	return fmt.Sprintf("Weather in %s:\nTemperature: %.2f°C\nCondition: %s",
		city, tempCelsius, weather.Weather[0].Description), nil
}

func sendDailyWeatherUpdates(bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		weatherInfo, err := getWeatherUpdate("Moscow")
		if err != nil {
			log.Printf("Error getting weather: %v", err)
			continue
		}

		for chatID := range chatIDs {
			msg := tgbotapi.NewMessage(chatID, weatherInfo)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending message to %d: %v", chatID, err)
			}
		}
	}
}

func handleMessages(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	userState := make(map[int64]string) // Добавляем мапу для отслеживания состояния пользователя

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		messageText := update.Message.Text

		switch {
		case messageText == "/start":
			chatIDs[chatID] = true
			userState[chatID] = "waiting_for_city"
			msg := tgbotapi.NewMessage(chatID, "Welcome! Enter your city name:")
			bot.Send(msg)

		case messageText == "/stop":
			delete(chatIDs, chatID)
			delete(userState, chatID)
			msg := tgbotapi.NewMessage(chatID, "You've unsubscribed from weather updates.")
			bot.Send(msg)

		case messageText == "/weather":
			msg := tgbotapi.NewMessage(chatID, "Please enter city name:")
			userState[chatID] = "waiting_for_city"
			bot.Send(msg)

		case userState[chatID] == "waiting_for_city":
			weatherInfo, err := getWeatherUpdate(messageText)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "City not found. Please try again.")
				bot.Send(msg)
				continue
			}
			msg := tgbotapi.NewMessage(chatID, weatherInfo)
			userState[chatID] = ""
			bot.Send(msg)
		}
	}
}

func main() {
	loadEnv()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	go sendDailyWeatherUpdates(bot)
	handleMessages(bot)
}
