# Registration
curl -X POST localhost:8080/api/user/register --data '{"login": "user7", "password": "123"}' -v

# Login
curl -X POST localhost:8080/api/user/login --data '{"login": "user1", "password": "123"}' -v

# Order Registration
curl -X POST localhost:8080/api/user/orders --data '123' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjExMTI3NDcsIlVzZXJJRCI6MSwiVXNlckxvZ2luIjoidXNlcjEifQ.HA11qV_7gKDaVTDdCnQPAJhrP1jHm3WVfnjVxYdNHcE'
