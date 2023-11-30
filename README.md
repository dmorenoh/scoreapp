# scoreapp

## Project setup
```docker-compose up -d```
Project will be available on localhost:8080

## Endpoints
### POST /user/{user_id}/score
Submits score for user

### GET /ranking/?type={type}
Returns ranking for given type. 
Type can be: "top100" or "At100/3" style