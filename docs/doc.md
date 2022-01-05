# AQI Server API
## AQI Station 
This API can be used to get/search for the station by many way
### AQI Station Get
```http request
GET /aqi/station
```
#### Params
| Field | Type   | Required        | Description                                                |
|-------|--------|-----------------|:-----------------------------------------------------------|
| qType | string | true            | The query type for request, must be "_get"                 |
| pType | string | true            | The query method, must be one of pType enum                |
| sid   | string | when pType=sid  | The station sid is number, from 0                          |
| name  | string | when pType=name | The full name or brief name of the station, like "beijing" |
| city  | string | when pType=city | The full name or brief name of a city, like "beijing"      |
| ip    | string | when pType=ip   | Can be ipv4/ipv6 like 192.168.1.1 or ipv6 address          |
| lon   | double | when pType=loc  | longitude, between -180 to 180 degree                      |
| lat   | double | when pType=loc  | latitude, between -90 to 90 degree                         |

#### PType Enum
| Value | Description                                |
|-------|--------------------------------------------|
| sid   | station unique identify                    |
| name  | station name                               |
| city  | the city where station locate at           |
| ip    | get station by ip  address                 |
| loc   | get station by longitude/latitude location |

#### Sample 
##### Request
```http request
GET http://127.0.0.1:9600/api/v1/aqi/station?qType=_get&pType=sid&sid=0
```
##### Response <font color=#2fc>OK</font>
```json lines
{
  "status": "OK",
  "code": 200,
  "body": {
    "sid": "0", // station identify
    "idx": 0, // station identify number
    "name": "Barrie, Ontario, Canada", // station name
    "loc": {
      "lon": -79.702306, // station longitude
      "lat": 44.382361  // station latitude
    }, 
    "up_time": 1641358800000, // station update time utc
    "tms": "2022-01-05T00:00:00-05:00", // station update time format
    "tz": "-05:00", // station timezone
    "city_name": "CA:Ontario/Barrie",  // station city name with country brief 
    "his_range": "2014-01~2021-12", 
    "sources": "[{\"name\":\"Citizen Weather Observer Program (CWOP/APRS)\",\"url\":\"http://wxqa.com/\",\"pols\":[\"weather\"],\"logo\":\"\"},{\"name\":\"Air Quality Ontario - the Ontario Ministry of the Environment and Climate Change\",\"url\":\"http://www.airqualityontario.com/\",\"pols\":null,\"logo\":\"Ontario-Ministry-of-the-Environment-and-Climate-Change.png\"}]"
  },
  "msg": "Success",
  "time": 1641360645662
}
```
##### Response <font color=#f22>ERROR</font>
```json
{
  "status": "Bad Request",
  "code": 400,
  "body": [
    {
      "FailedField": "StationGetRequest.Sid",
      "Tag": "number",
      "Value": "-1"
    }
  ],
  "msg": "",
  "time": 1641360974723
}
```
### AQI Station Search
```http request
GET /aqi/stations
```
#### Params
| Field | Type   | Required        | Description                                   |
|-------|--------|-----------------|:----------------------------------------------|
| qType | string | true            | The query type for request, must be "_search" |
| pType | string | true            | The query method, must be one of pType enum   |
## AQI Realtime

## AQI Forecast

## AQI History

## AQI LOGO