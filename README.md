# BookingService  
  
API Mircoservice for creating and booking meeting rooms for cross company events.

* Users can create and manage rooms
* Users can see meeting rooms availability  
* Users can book meeting rooms by the hour (first come first served)  
* Users can cancel their own reservations

**Requirements**:
* Go 1.15
* postgres
  
**Swagger**: [swagger.json](http://ec2-54-81-59-73.compute-1.amazonaws.com/api/apidocs/?url=http://ec2-54-81-59-73.compute-1.amazonaws.com/api/swagger.json)

## Linting & Testing
```
make lint
make test
```

##  Run
The service requires the export of `DATABASE_URL` at run time to properly configure Database connection.

### Local
``` 
make DATABASE_URL="postgres://postgres:test@127.0.0.1:5432/booking?sslmode=disable" \
 start local
```

### Docker
```
make docker-build
make DATABASE_URL"..." docker-start
```

### Examples

```
# Add Rooms
$ curl -X POST http://redfishbluefish.dev/rooms --data '{"Company":"coke","Number":3}' --header "Content-Type: application/json"
200 OK

# Get All Rooms
$ curl -X GET http://redfishbluefish.dev/rooms/all
[{"ID":1,"Name":"C1","Number":1,"Company":"C"},...]


# Get Room
$ curl -X GET http://redfishbluefish.dev/rooms/1
{"ID":1,"Name":"C1","Number":1,"Company":"C"}

# Delete Room
$ curl -X DELETE http://redfishbluefish.dev/rooms/1
200 OK

# Create Meeting
$ curl -X POST http://redfishbluefish.dev/booking --data '{"RoomID":1,"Title":"Meeting1", "Attendees":["alice","bob"],"Start":"2021-07-02T01:00:00Z"}' --header "Content-Type: application/json"
200 OK

# Get All Meetings
$ curl -X GET http://redfishbluefish.dev/booking/meetings/all
[{"ID":1,"RoomID":1,"Room":{"ID":1,"Name":"C1","Number":1,"Company":"C"},"Title":"Meeting1","Attendees":["alice","bob"],"Created":"2021-07-03T00:07:26.680792Z","Start":"2021-07-02T01:00:00Z","End":"2021-07-02T02:00:00Z"},...]

# Get Meeting
$ curl -X GET http://redfishbluefish.dev/booking/meetings/1
{"ID":1,"RoomID":1,"Room":{"ID":1,"Name":"C1","Number":1,"Company":"C"},"Title":"Meeting1","Attendees":["alice","bob"],"Created":"2021-07-03T00:07:26.680792Z","Start":"2021-07-02T01:00:00Z","End":"2021-07-02T02:00:00Z"}

# Delete Meeting
$ curl -X DELETE http://redfishbluefish.dev/booking/meetings/1
200 OK
```

### TODO:
* User models
* User Authentication
* Test pg-go and repositories
