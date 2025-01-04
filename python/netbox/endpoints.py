#!/bin/python3
import requests
# region TEST
def get_api_endpoints(url, headers):
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        return response.json()
    else:
        return None

def main():
    auth = f'Token 469f759c2e50ca1c998066fbd2ea360518ef7629'
    url = "http://cms01.mydomain.com:8000/api/"  # replace with your Netbox API URL
    headers = {"Authorization": auth}  # replace with your API token
    endpoints = get_api_endpoints(url, headers)
    if endpoints is not None:
        for key in endpoints.keys():
            print(f"Endpoint: {key}, URL: {endpoints[key]}")
    else:
        print("Failed to get API endpoints")

if __name__ == "__main__":
    main()


# endregion 