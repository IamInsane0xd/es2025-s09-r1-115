# Containers - Backend: Additionally added endpoints

### GET /api/containers/search
Can be used to filter through the containers using queries.<br>
<br>
If no query is provided, or a query contains malformed data (for example `?blockId=apple`), the endpoint will respond
with a `400 bad request`.<br>
If no container was found with the given criteria, the endpoint will respond with a `404 not found` error.<br>
If the query succeeded, the response will contain a list of containers matching the given
query.<br>
<br>
Queries that can be used:
  - blockId
  - bayNum
  - stackNum
  - tierNum
  - id[^1]

Example: 
```http request
###
GET http://localhost:3001/api/containers/search?bayNum=5&tierNum=5&blockId=1

HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 217

[
  {
    "id": "CFWY4364557",
    "blockId": 1,
    "bayNum": 5,
    "stackNum": 2,
    "tierNum": 4,
    "arrivedAt": "2024-01-20T00:00:00.000Z"
  },
  {
    "id": "COJW8239287",
    "blockId": 1,
    "bayNum": 5,
    "stackNum": 3,
    "tierNum": 4,
    "arrivedAt": "2024-01-17T00:00:00.000Z"
  }
]

###
GET http://localhost:3001/api/containers/search?blockId=1&id=UEKP3858709

HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 109

[
    {
        "id": "UEKP3858709",
        "blockId": 3,
        "bayNum": 2,
        "stackNum": 4,
        "tierNum": 1,
        "arrivedAt": "2024-01-06T00:00:00.000Z"
    }
]

###
GET http://localhost:3001/api/containers/search

HTTP/1.1 400 Bad Request
Content-Length: 15
Content-Type: text/plain; charset=utf-8

400 bad request
```

---

[^1]: If the id query is found, the others will be ignored
