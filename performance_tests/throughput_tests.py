import requests

url = "http://localhost:8080/trigger"

payload = "{\n\t\"gitrepo\": \"https://github.com/AartiJivrajani/Auto-Pub-Sub.git\"\n}"
headers = {
    'Content-Type': 'application/json'
}

response = requests.request("POST", url, headers=headers, data=payload)

print(response.text)
