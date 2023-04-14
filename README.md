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
### Historical percentages of renewables
## Webhook<br>
### Registration of webhook
### Deletion of webhook
### View registered webhook
## Status endpoint<br>
