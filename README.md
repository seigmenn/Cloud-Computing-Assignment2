
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


# General notes
Be mindful that the specified endpoint for this service may not be the exact same as the one used in the assignment description.

# Repository layout
```
.
├── Handler                    : Module for API
│   ├── CSVhandler.go          :
│   ├── consts.go              : Internal consts
│   ├── notifications.go       : Handler for weebhooks
│   ├── renewables.go          : Handler for renewables API
│   ├── status.go              : Handler for status API
│   └── structs.go             : Internal API structs
├── db_webhooks.go             :
├── go.mod                     :
├── main.go                    : Main application
├── renewable-share-energy.csv : Data used by handler
└── README.md
```
# Endpoints
## Renewable percentages<br>
### Current percentage of renewables
```
Path: /energy/v1/renewables/current/{country}?{neighbours=bool}
```
The endpoint returns the latest percentage of renewables in the energy mix.

This endpoint includes the parameters:<br>
`{country}` - a 3-letter country code (optional)<br>
`{neighbours=bool}` - a bool indicating whether the values of neighbouring countries should be shown (optional)<br>
<br>Example request with country code and neighbours:<br>
`/energy/v1/renewables/current/nor?neighbours=true`<br>

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
Number of matches: 1
Found in 0s
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
Path: energy/v1/renewables/history/{country}{begin=year&end=year}
```
This endpoint returns all the historical percentages of renewables in the energy mix. If no country is specified each country will return the mean value of all the historical percentages in the energy mix.


This endpoint includes the parameters:<br>
`{country}` - a 3-letter country code (optional)<br>
`{begin=year&end=year}` - limit the data returned to be within these two years (optional)<br>

Example request with country code:<br>
`energy/v1/renewables/history/usa`<br>

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
Number of results: 3
Found in 10ms
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
Number of matches: 103
Found in 0s
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
```
Method: POST
Path: /energy/v1/notifications/
```

This endpoints allows users to register webhooks that are triggered when a specified event is invoked. It is also possible to specify the number of invocations needed before a notification is triggered (i.e if the number is specified as 5, a notification will be called after 5,10,15 ... invocations)

Example message:
```

```

### Deletion of webhook
```
Method: DELETE
Path: /energy/v1/notifications/{id}
```

This endpoint deletes the webhook with the id `{id}`. <br>
Example message:
```

```

Example response:
```
```

### View registered webhook
```
Method: GET
Path: /energy/v1/notifications/{id}
```
This endpoint allows a user to view registered webhooks. If an id is specified it will only return that webhook, but if no id is specified it will return all registered endpoints. <br>

Example message with an id:
```
```

Example response:
```
```

Example message without id:
```
```

Example response:
```
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

```


