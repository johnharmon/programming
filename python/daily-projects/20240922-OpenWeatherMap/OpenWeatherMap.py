#!/bin/python3

import sys 
import requests
import json 
import yaml
import os


class locationInfo():
    def __init__(self,city, state, country):
        self.city = city 
        self.state = state
        self.country = country 

def getLatLon(city, apiKey = ''):

    apiUrl = f'http://api.openweathermap.org/geo/1.0/direct?q={city}&limit=1&appid={apiKey}'
    print(apiUrl)
    response = requests.get(url = apiUrl, headers = {'accept': 'application/json'})
    responseJson = response.json() 
    result = responseJson[0]
    return {'lat': result['lat'], 'lon': result['lon'], 'country': result['country'], 'state': result['state'] }
    #qString = '.'.join(location)
    #http://api.openweathermap.org/geo/1.0/direct?q=London&limit=5&appid={API key}
    #print(response)
    #print(json.dumps(response.json(), indent=2))

def queryCity(city, apiKey):
    cityInfo = getLatLon(city, apiKey)
    cityUrl = f'https://api.openweathermap.org/data/2.5/weather?lat={cityInfo["lat"]}&lon={cityInfo["lon"]}&appid={apiKey}'
    response = requests.get(url = cityUrl, headers = {'accept': 'applicaiton/json'})
    responseJson = response.json()
    responseText = response.text
    #print(yaml.dump(response.json(), default_flow_style=False))
    displayCity(responseJson)

def notEmpty(s):
    if len(s.strip()) == 0:
        return False 
    else: 
        return True

def displayCity(cityData):
    print(f'{cityData["name"]}, {cityData["sys"]["country"]} currently:')
    print(f'  Feels like: {round(cityData["main"]["feels_like"] - 271, 2)} C')
    print(f'  Actual Temp: {round(cityData["main"]["temp"] - 271, 2)} C')
    print(f'  High: {round(cityData["main"]["temp_max"] - 271, 2)} C')
    print(f'  Low: {round(cityData["main"]["temp_min"] - 271, 2)} C')
    print(f'  Humidity: {cityData["main"]["humidity"]}%')
    print(f'  Weather: {cityData["weather"][0]["description"]}')
    print(f'  Wind:')
    print(f'    Deg: {cityData["wind"]["deg"]},  Speed: {cityData["wind"]["speed"]} kts')

def main():
    city = sys.argv[1]
    apiKeyFile = input("Please provide the path to the api key file: ") 
    apiKeyFile = os.path.expandvars(os.path.expanduser(apiKeyFile))
    with open(apiKeyFile, 'r') as akf:
        apiKey = akf.read().strip()
    queryCity(city, apiKey)

    #location = {}
    #location['city'] = input("Please enter the city name: ")
    #location['state'] = input("Please enter the state code (US only) (Optional): ")
    #location['country'] = input("Please enter the country code (ISO 3166 codes only) (Optional): ")
    #suppliedValues = [key for key in filter(lambda x: len(x.strip()) > 0, location)]

if __name__ == '__main__':
    main()
