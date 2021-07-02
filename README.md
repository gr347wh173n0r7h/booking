# BookingService  
  
API Mircoservice for creating and booking meeting rooms for cross company events.

* Users can create and manage rooms
* Users can see meeting rooms availability  
* Users can book meeting rooms by the hour (first come first served)  
* Users can cancel their own reservations

**Requirements**:
* Go 1.15
* postgres
  
**Swagger**: TBD

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

### TODO:
* User models
* User Authentication
* Test pg-go and repositories
