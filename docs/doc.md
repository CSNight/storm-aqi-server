# AQI Server API
## Common Enum
### Pollutant Enum
| Value | Description                      |
|-------|----------------------------------|
| co    | carbon monoxide                  |
| no2   | Nitrogen dioxide                 |
| pm25  | Pulmonary particulate matter 2.5 |
| pm10  | Pulmonary particulate matter 10  |
| o3    | ozone                            |
| so2   | sulfur dioxide                   |

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
| pType | string | true            | The query method. See pType enum                           |
| sid   | string | when pType=sid  | The station sequence id number, from 0                     |
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
GET http://aqiserver/api/v1/aqi/station?qType=_get&pType=sid&sid=0
```
##### Response 200 <font color=#2f5>OK</font>
```json lines
{
  "status": "OK",
  "code": 200,
  "body": {
    "sid": "0", // the station sequence id number
    "idx": 0, // the station sequence id number
    "name": "Barrie, Ontario, Canada", // the station name
    "loc": {  // the station location coordinate
      "lon": -79.702306,
      "lat": 44.382361
    },
    "up_time": 1641445200000, // the station update timestamp in utc
    "tms": "2022-01-06T00:00:00-05:00", // the station update time in rfc2822 format
    "tz": "-05:00", // the station update timezone
    "city_name": "CA:Ontario/Barrie", // the city where station belongs
    "his_range": "2014-01~2021-12", 
    "sources": [ // the station metadata
      {
        "name": "Citizen Weather Observer Program (CWOP/APRS)",
        "url": "http://wxqa.com/",
        "pols": [
          "weather"
        ]
      },
      {
        "logo": "Ontario-Ministry-of-the-Environment-and-Climate-Change.png",
        "name": "Air Quality Ontario - the Ontario Ministry of the Environment and Climate Change",
        "url": "http://www.airqualityontario.com/"
      }
    ]
  },
  "msg": "Success",
  "time": 1641448426986
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
| Field       | Type     | Required          | Description                                                 |
|-------------|----------|-------------------|:------------------------------------------------------------|
| qType       | string   | true              | The query type for request, must be "_search"               |
| pType       | string   | true              | The query method. See pType enum                            |
| size        | int      | true              | size of stations in response, maximum support 10000         |
| name        | string   | when pType=name   | The full name or brief name of the station, like "beijing"  |
| city        | string   | when pType=city   | The full name or brief name of a city, like "beijing"       |
| topLeft     | []double | when pType=area   | The bound top left corner lon/lat coordinate like 80,39     |
| bottomRight | []double | when pType=area   | The bound bottom right corner lon/lat coordinate like 80,39 |
| center      | []double | when pType=radius | The center of cycle area lon/lat coordinate like 80,39      |
| radius      | double   | when pType=radius | The radius of cycle area to search, maximum support 10000   |
| unit        | string   | when pType=radius | The unit of radius, must be one of kilometers miles meters  |
#### PType Enum
| Value  | Description                 |
|--------|-----------------------------|
| name   | search stations name        |
| city   | search stations by city     |
| area   | search stations by envelope |
| radius | search stations by cycle    |
#### Sample
##### Request
```http request
GET http://aqiserver/api/v1/aqi/stations?qType=_search&pType=name&size=1&name=beijing
```
##### Response 200 <font color=#2f5>OK</font>
```json
{
  "status": "OK",
  "code": 200,
  "body": [
    {
      "sid": "1451",
      "idx": 1451,
      "name": "Beijing (北京)",
      "loc": {
        "lon": 116.468117,
        "lat": 39.954592
      },
      "up_time": 1641445200000,
      "tms": "2022-01-06T13:00:00+08:00",
      "tz": "+08:00",
      "city_name": "Beijing",
      "his_range": "2014-01~2021-12",
      "sources": [
        {
          "name": "Citizen Weather Observer Program (CWOP/APRS)",
          "url": "http://wxqa.com/",
          "pols": [
            "weather"
          ]
        },
        {
          "name": "Beijing Environmental Protection Monitoring Center (北京市环境保护监测中心)",
          "url": "http://www.bjmemc.com.cn/"
        }
      ]
    }
  ],
  "msg": "Success",
  "time": 1641453247827
}
```

## AQI Realtime

### AQI Realtime Get
```http request
GET /aqi/realtime
```
#### Params
| Field | Type   | Required          | Description                                                                          |
|-------|--------|-------------------|:-------------------------------------------------------------------------------------|
| qType | string | true              | The query type for request, must be "_get"                                           |
| pType | string | true              | The query method, must be one of all/single means all pollutants or single pollutant |
| sid   | string | true              | The station sequence id number, from 0                                               |
| pol   | string | when pType=single | The pollutant type want to get. See Pollutant Enum                                   |
#### Sample
##### Request
```http request
GET http://aqiserver/api/v1/aqi/realtime?qType=_get&pType=all&sid=0
```
##### Response 200 <font color=#2f5>OK</font>
```json lines
{
  "status": "OK",
  "code": 200,
  "body": {
    "idx": 0,
    "sid": "0",
    "name": "Barrie, Ontario, Canada",
    "loc": {
      "lon": -79.702306,
      "lat": 44.382361
    },
    "city_name": "CA:Ontario/Barrie",
    "realtime": [
      {
        "pol": "o3", // pollutant type
        "data": 17.6 // value
      },
      {
        "pol": "no2",
        "data": 3.4
      },
      {
        "pol": "pm25",
        "data": 25
      },
      {
        "pol": "so2",
        "data": 0.2
      }
    ],
    "tz": "-05:00",
    "tm": 1641445200000, //  pollutant value last update timestamp in utc
    "tms": "2022-01-06T00:00:00-05:00" // last update time in rfc2822 format
  },
  "msg": "Success",
  "time": 1641455505471
}
```
## AQI Forecast
### AQI Forecast Get
```http request
GET /aqi/forecast
```
#### Params
| Field | Type   | Required          | Description                                                                          |
|-------|--------|-------------------|:-------------------------------------------------------------------------------------|
| qType | string | true              | The query type for request, must be "_get"                                           |
| pType | string | true              | The query method, must be one of all/single means all pollutants or single pollutant |
| sid   | string | true              | The station sequence id number, from 0                                               |
| pol   | string | when pType=single | The pollutant type want to get. See Pollutant Enum                                   |
#### Sample
##### Request
```http request
GET http://aqiserver/api/v1/aqi/forecast?qType=_get&pType=single&sid=0&pol=pm25
```
##### Response 200 <font color=#2f5>OK</font>
```json lines
{
  "status": "OK",
  "code": 200,
  "body": {
    "idx": 0,
    "sid": "0",
    "name": "Barrie, Ontario, Canada",
    "loc": {
      "lon": -79.702306,
      "lat": 44.382361
    },
    "city_name": "CA:Ontario/Barrie",
    "forecast": { // aqi forecast value mapped by pollutant
      "pm25": [
        {
          "avg": 35, // average value
          "day": "2022-01-04", // forecast day
          "max": 40, // maximum value
          "min": 19 // minimum value
        },
        {
          "avg": 19,
          "day": "2022-01-05",
          "max": 45,
          "min": 3
        },
        {
          "avg": 3,
          "day": "2022-01-06",
          "max": 3,
          "min": 3
        },
        {
          "avg": 3,
          "day": "2022-01-07",
          "max": 3,
          "min": 3
        },
        {
          "avg": 12,
          "day": "2022-01-08",
          "max": 33,
          "min": 4
        },
        {
          "avg": 9,
          "day": "2022-01-09",
          "max": 17,
          "min": 3
        },
        {
          "avg": 3,
          "day": "2022-01-10",
          "max": 3,
          "min": 3
        }
      ]
    },
    "tz": "-05:00",
    "tm": 1641452400000,
    "tms": "2022-01-06T02:00:00-05:00"
  },
  "msg": "Success",
  "time": 1641457408192
}
```
## AQI History
### AQI History Get
```http request
GET /aqi/history
```
#### Params
| Field  | Type   | Required          | Description                                                                                                 |
|--------|--------|-------------------|:------------------------------------------------------------------------------------------------------------|
| qType  | string | true              | The query type for request, must be "_get"                                                                  |
| pType  | string | true              | The query method, must be one of recent range                                                               |
| sid    | string | true              | The station sequence id number, from 0                                                                      |
| pol    | string | true              | The pollutant type want to get. must be Pollutant Enum or "all" for all pollutant                           |
| recent | string | when pType=recent | The recent time range. See Recent Enum                                                                      |
| start  | string | when pType=range  | The custom time range start point with format like 2022-01-01 <br/>search will include endpoint of therange |
| end    | string | when pType=range  | The custom time range end point with format like 2022-01-01                                                 |
#### PType Enum
| Value  | Description                                                  |
|--------|--------------------------------------------------------------|
| recent | Search the station history by a shortcut time range from now |
| range  | Search the station by custom time range                      |
#### Recent Enum
| Value       | Description      |
|-------------|------------------|
| lastDay     | yesterday        |
| lastWeek    | last week        |
| lastMonth   | last month       |
| lastQuarter | last three month |
| lastYear    | last year        |
#### Sample by Range
##### Request 
```http request
GET http://aqiserver/api/v1/aqi/history?qType=_get&pType=range&sid=0&start=2021-09-02&end=2021-09-03&pol=all
```
##### Response 200 <font color=#2f5>OK</font>
```json lines
{
  "status": "OK",
  "code": 200,
  "body": {
    "idx": 0,
    "sid": "0",
    "name": "Barrie, Ontario, Canada",
    "loc": {
      "lon": -79.702306,
      "lat": 44.382361
    },
    "city_name": "CA:Ontario/Barrie",
    "history": { // station history info mapped by pollutant
      "co": [], // not found history for specification pollutant and time range
      "no2": [
        {
          "pol": "no2", // pollutant type
          "name": "NO<sub>2</sub>", // pollutant name with subscript
          "data": 3, //value
          "tz": "-5.00", // timezone
          "month": 9, // history in month
          "year": 2021, // history in year
          "tm": 1630627200000, // history timestamp in utc
          "tms": "2021-09-02T19:00:00-05:00" // history time in rfc2822 format
        },
        {
          "pol": "no2",
          "name": "NO<sub>2</sub>",
          "data": 2,
          "tz": "-5.00",
          "month": 9,
          "year": 2021,
          "tm": 1630540800000,
          "tms": "2021-09-01T19:00:00-05:00"
        }
      ],
      "o3": [
        {
          "pol": "o3",
          "name": "O<sub>3</sub>",
          "data": 13,
          "tz": "-5.00",
          "month": 9,
          "year": 2021,
          "tm": 1630627200000,
          "tms": "2021-09-02T19:00:00-05:00"
        },
        {
          "pol": "o3",
          "name": "O<sub>3</sub>",
          "data": 18,
          "tz": "-5.00",
          "month": 9,
          "year": 2021,
          "tm": 1630540800000,
          "tms": "2021-09-01T19:00:00-05:00"
        }
      ],
      "pm10": [],
      "pm25": [
        {
          "pol": "pm25",
          "name": "PM<sub>2.5</sub>",
          "data": 15,
          "tz": "-5.00",
          "month": 9,
          "year": 2021,
          "tm": 1630627200000,
          "tms": "2021-09-02T19:00:00-05:00"
        },
        {
          "pol": "pm25",
          "name": "PM<sub>2.5</sub>",
          "data": 24,
          "tz": "-5.00",
          "month": 9,
          "year": 2021,
          "tm": 1630540800000,
          "tms": "2021-09-01T19:00:00-05:00"
        }
      ],
      "so2": [
        {
          "pol": "so2",
          "name": "SO<sub>2</sub>",
          "data": 0,
          "tz": "-5.00",
          "month": 9,
          "year": 2021,
          "tm": 1630627200000,
          "tms": "2021-09-02T19:00:00-05:00"
        }
      ]
    }
  },
  "msg": "Success",
  "time": 1641458411235
}
```
#### Sample by Range
##### Request
```http request
GET http://aqiserver/api/v1/aqi/history?qType=_get&pType=range&sid=0&start=2021-09-02&end=2021-09-03&pol=all
```
##### Response 200 <font color=#2f5>OK</font>
```json lines
{
  "status": "OK",
  "code": 200,
  "body": {
    "idx": 0,
    "sid": "0",
    "name": "Barrie, Ontario, Canada",
    "loc": {
      "lon": -79.702306,
      "lat": 44.382361
    },
    "city_name": "CA:Ontario/Barrie",
    "history": {
      "co": [],
      "no2": [
        {
          "pol": "no2",
          "name": "NO<sub>2</sub>",
          "data": 9.02,
          "tz": "-05:00",
          "month": 1,
          "year": 2022,
          "tm": 1641340800000,
          "tms": "2022-01-04T19:00:00-05:00"
        }
      ],
      "o3": [
        {
          "pol": "o3",
          "name": "O<sub>3</sub>",
          "data": 14.78,
          "tz": "-05:00",
          "month": 1,
          "year": 2022,
          "tm": 1641340800000,
          "tms": "2022-01-04T19:00:00-05:00"
        }
      ],
      "pm10": [],
      "pm25": [
        {
          "pol": "pm25",
          "name": "PM<sub>2.5</sub>",
          "data": 26.36,
          "tz": "-05:00",
          "month": 1,
          "year": 2022,
          "tm": 1641340800000,
          "tms": "2022-01-04T19:00:00-05:00"
        }
      ],
      "so2": []
    }
  },
  "msg": "Success",
  "time": 1641458815507
}
```
## AQI Logo
### AQI Station Logo Get
```http request
GET /aqi/logo/{logoName}
```
#### Params
| Field    | Type   | Required | Description                      |
|----------|--------|----------|:---------------------------------|
| logoName | string | true     | logoName in station source field |
#### Sample
##### Request
```http request
GET http://aqiserver/api/v1/aqi/logo/Ontario-Ministry-of-the-Environment-and-Climate-Change.png
```