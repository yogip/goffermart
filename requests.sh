# Registration
curl -X POST localhost:8080/api/user/register --data '{"login": "user7", "password": "123"}' -v

# Login
curl -X POST localhost:8080/api/user/login --data '{"login": "user1", "password": "123"}' -v

# Order Registration
curl -X POST localhost:8080/api/user/orders --data '123' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjExMTI3NDcsIlVzZXJJRCI6MSwiVXNlckxvZ2luIjoidXNlcjEifQ.HA11qV_7gKDaVTDdCnQPAJhrP1jHm3WVfnjVxYdNHcE'

# Get balance
curl localhost:8080/api/user/balance -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE1NzEwNjIsIlVzZXJJRCI6MSwiVXNlckxvZ2luIjoidXNlcjEifQ.PFzPIlmunhyBL9pgkLZDOXEBDYBRsC3Lc8b5mLDuE_A'

# Withdrawn list
curl -X POST localhost:8080/api/user/balance/withdraw --data '{"order": 23772524, "sum": 1}' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE1NzEwNjIsIlVzZXJJRCI6MSwiVXNlckxvZ2luIjoidXNlcjEifQ.PFzPIlmunhyBL9pgkLZDOXEBDYBRsC3Lc8b5mLDuE_A'

# Withdrawn
curl -X GET localhost:8080/api/user/withdrawals -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE1NzEwNjIsIlVzZXJJRCI6MSwiVXNlckxvZ2luIjoidXNlcjEifQ.PFzPIlmunhyBL9pgkLZDOXEBDYBRsC3Lc8b5mLDuE_A'