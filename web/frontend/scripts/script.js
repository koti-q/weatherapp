// DOM Elements
const cityInput = document.getElementById('city-input');
const searchButton = document.getElementById('search-button');

const weatherInfo = document.getElementById('weather-info');
const cityName = document.getElementById('city-name');
const temperature = document.getElementById('temperature');
const humidity = document.getElementById('humidity');
const windSpeed = document.getElementById('wind-speed');

// Error message elements
const errorMessage = document.getElementById('error-message');
const errorMessageText = document.getElementById('error-message-text');

// Fetch weather data from the API
async function fetchWeather(city) {
    const API_URL = `http://localhost:8080/api/weather?city=${city}`;

    try {
        const response = await fetch(API_URL);

        if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const weatherData = await response.json();
        console.log('Weather Data:', weatherData);
        return weatherData;
    } catch (error) {
        console.error('Error fetching weather:', error.message);
        return null;
    }
}

// Display weather data on the page
function displayWeather(weatherData) {
    if (!weatherData) {
        errorMessage.style.display = 'block';
        errorMessageText.textContent = 'Error fetching weather data. Please try again.'; 
        console.error('No weather data to display.');
        return;
    }
    temp = Math.round(weatherData.main.temp - 273.15);

    // Update DOM elements with weather data
    cityName.textContent = weatherData.name;
    temperature.textContent = `Temperature: ${temp}°C`;
    humidity.textContent = `Humidity: ${weatherData.main.humidity}%`;
    windSpeed.textContent = `Wind Speed: ${weatherData.wind.speed} m/s`;

    // Show the weather info section
    weatherInfo.style.display = 'block';
}

weatherInfo.style.display = 'none';

// Event listener for the search button
searchButton.addEventListener('click', async () => {
// Clear any previous error messages
    errorMessage.style.display = 'none';
    const city = cityInput.value.trim();

    if (!city) {
    // Display an error message
        errorMessage.style.display = 'block';
        errorMessageText.textContent = 'Please enter a city name.';

        console.error('Please enter a city name.');
        return;
    }

    console.log('Searching for weather in:', city);

    try {
        const weatherData = await fetchWeather(city);

        if (weatherData) {
            displayWeather(weatherData);
        }
    } catch (error) {
        errorMessage.style.display = 'block';
        errorMessageText.textContent = 'An error occurred while fetching weather data.';
        console.error('Error:', error);
    }
});
