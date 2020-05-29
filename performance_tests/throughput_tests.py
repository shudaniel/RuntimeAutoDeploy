# importing the requests library 
import requests 
  
# api-endpoint 
URL = "localhost:8080/trigger"
git_repo = "https://github.com/AartiJivrajani/Auto-Pub-Sub.git"  
  
# defining a params dict for the parameters to be sent to the API 
PARAMS = {'gitrepo': git_repo} 
  
# sending get request and saving the response as response object 
r = requests.post(url = URL, params = PARAMS) 
  
# extracting data in json format 
data = r.test
  
print(data)