package Handler

import "errors"

const LINEBREAK = "\n"                                              //Line break
const BASEPATH = "/energy/v1/"                                      //Base path of webservice
const RENEWABLESPATH = BASEPATH + "renewables/"                     //Path to renewables endpoint
const NOTIFICATIONSPATH = BASEPATH + "notifications/"               //Path to notifications endpoint
const STATUSPATH = BASEPATH + "status"                              //Path to status endpoint
const COUNTRIESAPIALPHA = "http://129.241.150.113:8080/v3.1/alpha/" //URL of Countries API for search with alpha-code
const COUNTRIESAPIURL = "http://129.241.150.113:8080"               //URL to root of Countries API for status endpoint
const CSVPATH = "renewable-share-energy.csv"
const COLLECTION = "Wirehooks" //Collection name in Firestore

var ERRFILEREAD = errors.New("couldn't read file")
var ERRFILEPARSE = errors.New("file couldn't be parsed")
var ERRCOUNTRYNOTFOUND = errors.New("no country found")
