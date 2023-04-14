# Overview of read.me
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
# Repository layout
```
.
├── Handler                    : Module for API
│   ├── CSVhandler.go          :
│   ├── consts.go              :
│   ├── notifications.go       :
│   ├── renewables.go          :
│   ├── status.go              :
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
```

Example request with country code and begin and end year:<br>
`energy/v1/renewables/history/nor`<br>

Example response:
```
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
### Deletion of webhook
### View registered webhook
## Status endpoint<br>
