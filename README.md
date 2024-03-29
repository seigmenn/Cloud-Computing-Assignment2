
# Overview of README.md
- General notes
- Repository layout
- Endpoints
	- Renewable percentages:
		- Current percentage of renewables
		- Historical percentages of renewables
	- Webhook:
		- Registration of webhook
		- Deletion of webhook
		- View registered webhook
	- Status endpoint
- Openstack and Docker deployment
	- Dockerfile

# General notes
The version of this app hosted on docker can be found here: http://10.212.172.23:8080/ (remember to be on the NTNU network)<br>
Be mindful that the specified endpoint for this service may not be the exact same as the one used in the assignment description.

# Repository layout
```
.
├── Handler                    : Module for API
│   ├── CSVhandler.go          : Handler for CSV-files
│   ├── consts.go              : Internal consts
│   ├── notifications.go       : Handler for weebhooks
│   ├── renewables.go          : Handler for renewables API
│   ├── status.go              : Handler for status API
│   └── structs.go             : Internal API structs
├── db_webhooks.go             : Modules
├── go.mod                     : Modules
├── main.go                    : Main application
├── renewable-share-energy.csv : Data used by handler
└── README.md
```
# Endpoints
## Renewable percentages<br>
### Current percentage of renewables
```
Path: /energy/v1/renewables/current/{country}?{neighbours=bool}{info=1}
```
The endpoint returns the latest percentage of renewables in the energy mix.

This endpoint includes the parameters:<br>
`{country}` - a 3-letter country code or a full country name (optional)<br>
`{neighbours=bool}` - a bool indicating whether the values of neighbouring countries should be shown (optional)<br>
`{info=1}` - if specified, number of matches and the time to return data(optional)<br>
<br>Example request with country code and neighbours:<br>
`/energy/v1/renewables/current/nor?neighbours=true&info=1`<br>

Example response:
```
Number of matches: 4
Found in 1.9s
[
 {
  "name": "Norway",
  "isoCode": "NOR",
  "year": 2021,
  "percentage": 71.55836486816406
 },
 {
  "name": "Finland",
  "isoCode": "FIN",
  "year": 2021,
  "percentage": 34.611289978027344
 },
 {
  "name": "Sweden",
  "isoCode": "SWE",
  "year": 2021,
  "percentage": 50.924007415771484
 },
 {
  "name": "Russia",
  "isoCode": "RUS",
  "year": 2021,
  "percentage": 6.620289325714111
 }
]

```

Example request with country code:<br>
`energy/v1/renewables/current/usa`<br>

Example response:
```

[
 {
  "name": "United States",
  "isoCode": "USA",
  "year": 2021,
  "percentage": 10.655990600585938
 }
]
```

### Historical percentages of renewables
```
Path: energy/v1/renewables/history/{country}?{begin=year&end=year}{info=1}
```
This endpoint returns all the historical percentages of renewables in the energy mix. If no country is specified each country will return the mean value of all the historical percentages in the energy mix.


This endpoint includes the parameters:<br>
`{country}` - a 3-letter country code (optional)<br>
`{begin=year&end=year}` - limit the data returned to be within these two years (optional)<br>
`{info=1}` - if specified, number of matches and the time to return data(optional)<br>

Example request with country code:<br>
`energy/v1/renewables/history/usa?info=1`<br>

Example response:
```
Number of results: 57
Found in 10ms
[
	 {
	  "name": "United States",
	  "isoCode": "USA",
	  "year": 1965,
	  "percentage": 4.368869781494141
	 },
	 {
	  "name": "United States",
	  "isoCode": "USA",
	  "year": 1966,
	  "percentage": 4.1714019775390625
	 },
	 {
	  "name": "United States",
	  "isoCode": "USA",
	  "year": 1967,
	  "percentage": 4.542215824127197
	 },
	 {
	  "name": "United States",
	  "isoCode": "USA",
	  "year": 1968,
	  "percentage": 4.3309736251831055
	 },
 ...
 ]
```

Example request with country code and begin and end year:<br>
`/energy/v1/renewables/history/nor?begin=1972&end=1974`<br>

Example response:
```
[
 {
  "name": "Norway",
  "isoCode": "NOR",
  "year": 1972,
  "percentage": 64.29580688476562
 },
 {
  "name": "Norway",
  "isoCode": "NOR",
  "year": 1973,
  "percentage": 65.58218383789062
 },
 {
  "name": "Norway",
  "isoCode": "NOR",
  "year": 1974,
  "percentage": 68.31411743164062
 }
]
```

Example request with no optional parameters (return mean values):<br>
`energy/v1/renewables/history`<br>

Example response:
```
[
 ...
	{
	  "name": "Australia",
	  "isoCode": "AUS",
	  "percentage": 5.300048171428212
	 },
	 {
	  "name": "Austria",
	  "isoCode": "AUT",
	  "percentage": 29.462373633133737
	 },
	 {
	  "name": "Azerbaijan",
	  "isoCode": "AZE",
	  "percentage": 3.2902767239390194
	 },
	 {
	  "name": "Bangladesh",
	  "isoCode": "BGD",
	  "percentage": 2.5659469354386424
	 }
 ...
 ]
```


## Webhook<br>

### Registration of webhook

Users, through the usage of the `/energy/v1/notifications/` endpoint, have the ability to do a number of actions regarding webhooks, such as the matter of registration. To register a webhook, the endpoint requires a `HTTP POST` request, with a body containing two/three values:

- `URL`; string - the URL the service hooks onto, and sends a request when invocation triggers.
- `CALLS`: integer - the amount of events that has to happen per before an invocation triggers (calls = 5 will return an invocation for every 5*n calls (5, 10, 15, ...)).
- `COUNTRY`: optional string - a three-letter long isocode which represents a country a webhook will be hooked onto. This is an optional value, and not including it will make the webhook apply for all countries.

Given a successful registration, the service will return a `application/json` of the webhook's ID, generated from hashing the URL of the hooked site. It is recommended that this ID gets saved for further usage of the program.

The registration of these webhooks have been stored on a remote database which in our case is Firebase. This is to ensure that the webhooks that have been registered will survive a service restart. 
On firebase, the database can have different collections, but we only have one collection for our webhooks. 
The collection consists of documents with uniquely autogenerated ID's and the data in the documents is the information the user has given when registering a webhook.

#### Request Example

```
Method: POST
Path: /energy/v1/notifications/
```
```
{
    "url": "localhost:8080/client/",
    "calls": 5,
    "country": "NOR"
}
```

#### Response Example
```
{
	"webhook_id": "8d5f66188edc3dd36776c02ec61632edcb677251939e002a0735204ffc25976e"
}
```

### Deletion of webhook

Users can also choose to delete registered webhooks, given they have the ID of a webhook they would want to delete. To do so, the user first goes to the endpoint of `/energy/v1/notifications/{id}`, where "id" is their webhook id. To this URL, the user sends a `HTTP DELETE` request, which will either return a `text/html` message of success (deleting the webhook), or an `error` (not deleting any webhooks). 

#### Request Example

```
Method: DELETE
Path: /energy/v1/notifications/{id}
```

#### Response Example

Successful response:
```
You are now deleting the following information:
Identification: {id}
URL: (webhook.URL)
In which it was focused to look at the country of (webhook.ISO) and report a notification every (webhook.Calls) invocations.
```


### View registered webhook

Users can also read information about all registered webhooks, either specifically mentioned ones through the usages of webhook_id, or all that are registered in the program.
To retrieve information about all webhooks, the user has to send a `HTTP GET` request to endpoint of `/energy/v1/notifications/`. To retrieve information about a specific one, the user has to add the ID of the webhook to the endpoint previously mentioned, and perform a `HTTP GET` request to that; `/energy/v1/notifications/{id}`.

The information retrieved is in the form of `application/json`, containing a webhook's URL, Calls and 

#### Example Requests

Retrieve information about all webhooks:

```
Method: GET
Path: /energy/v1/notifications/
```

Retrieve information about one specific webhook by ID:

```
Method: GET
Path: /energy/v1/notifications/{id}
```

#### Example Responses

Retrieve information about all webhooks:
```
[
    {
        "url": "https://localhost:8080/client/",
        "country": "NOR",
        "calls": 5,
        "webhook_id": "8d5f66188edc3dd36776c02ec61632edcb677251939e002a0735204ffc25976e"
    },
    {
        "url": "https://localhost:8080/client/",
        "country": "NOR",
        "calls": 10,
        "webhook_id": "1ec08701e65b7a91c7b34e06f9bcefe5e22b7657b9fd6fa0acff5f189f9811a1"
    }, ... 
]
```

Retrieve information about one specific webhook:
```
{
    "url": "https://localhost:8080/client/",
    "country": "NOR",
    "calls": 5,
    "webhook_id": "1ec08701e65b7a91c7b34e06f9bcefe5e22b7657b9fd6fa0acff5f189f9811a1"
}
```

### Webhook Invocation

When a webhook has been triggered, it will send a `HTTP POST` request to the specified URL belonging to the webhook - this request will then contain in the form of `application/json`; the webhook's ID, the full name of the country the webhook is registered to, and the amount of calls that the webhook has had up til that moment.

#### Example Response
```
Method: POST
Path: webhook.URL
```

```
{
	"webhook_id": "1ec08701e65b7a91c7b34e06f9bcefe5e22b7657b9fd6fa0acff5f189f9811a1",
	"country": "Norway",
	"calls": 10
}
```


## Status endpoint<br>
```
Path: /energy/v1/status/
```

This endpoint provides information about the availability of all the individual services that this service depends upon.

Things checked in the status endpoint include:<br>
- The notification database API
- The REST countries API
- Number of registered webhooks
- Version of this service
- Uptime of this service

Example response:
```
{
	countries_api: 200 - OK
	notification_db: 200 - OK
	webhooks: 6
	version: v1
	uptime: 17h5m44.4s
}

```

# Openstack and Docker deployment 
This assignment is hosted via Openstack and by the use of Docker. <br>
For this assignment, this git repo has been copied to a virtual machine hosted on Openstack. <br> Using a Dockerfile (see next point for further info) a docker image is created using the command `docker build --tag app:1 .` where `tag` sets the name `app` for the build and specifies the version of the app `1` and the `.` represents where the Dockerfile is located. <br> To run the image  you use the docker run command `docker run -d -p 8080:8080 --restart always app:1` where `-d` starts it in detached mode, `-p` opens up the specified port numbers, `--restart always` makes the image restart automatically and both the name of the image and the version is specified. 

## Dockerfile
The contents of the Dockerfile used to build the go-app:

```                                                                                  
FROM golang:1.19		#Specifying which version of go to use
WORKDIR /build			#Creating a work directory for the go build
ADD go.mod .			#Adds the external modules into the work directory
COPY . .
RUN go build -o /main 	#building the go app into a binary called main
EXPOSE 80				#Exposing the proper ports to make it accesible ecternally
CMD ["/main"]			#Setting the command that will be run when the image is run
```
