# poc-csv-persistence
Repository for a microservice responsible to do CRUD operations at a database.  
This repository is nothing to be really used in production. Is something just to demonstrate my Golang Skills
This intended to accept large files. used with [this repository](https://github.com/VictorPrado99/poc-csv-uploader), is possible to send a csv file, and wait the processing happen. The process is asynchronous.  

### Resource Definition

```json
[
    {
        "id": 1074315315,
        "email": "test1e@teste",
        "phone_number": "351 961 251 326",
        "parcel_weight": 22.4,
        "date": "2022-03-12",
        "country": "Portugal"
    }
]
```

### `GET` /orders

Get a paginated arrays of orders

**QueryParameters**

- *sort*
  This parameter will sort your content based on ID
  **default: asc
  domain: [asc, desc]**

- *offset*
This parameter will have the offset value from pagination purposes. The API will send you back a header field, with the next value. This header will be always, even if is the last offset chunk, the number will not change though.
**default: 0**

- *limit*
This parameter will set the limit for pagination
**default: 10**
**limit: 100**

- *country*
This parameter, will accept more than one value, separated by commas. You can filter by country your result

- *date*
Filter by date, you can filter justo for a single date

- *weightLimit*
Will retrieve everything of weight equal or less

### `POST` /orders

This endpoint will receive a array of order in the body, and will persist everything in the database

### `HEAD` /orders/{country}

Will retrieve the sum of weight and how many orders you have per country

e.g
**Header**
```json
{
    "x-country": "Portugal",
    "X-Orders-Count": 15,
    "X-Weight-Sum": 20.12
}
```


## Architeture

The server will look for a file at your ***$HOME*** named config.json, this file has the following pattern

```json
{
    "connection_string": "root:root@tcp(mysql:3306)/csv_persistence",
    "port": 9001
}
```

This place can be overwriten with the enviroment variable ***$PERSISTENCE_HOME*** 

The Docker Image is [available here](https://hub.docker.com/repository/docker/victorprado99/poc-csv-persistence)
 
### Docker-Compose

The ./dev directory, we have a docker-compose.yml intended to developing and unit tests. The file will run only a database instance which the api will connect.  
At the root dir, have docker-compose file which can run the entire application, just run

```bash
docker-compose up
```

And you are ready to make requests at the api, with postman, with [poc-persistence-uploader](https://github.com/VictorPrado99/poc-csv-uploader) or whatever client you like.