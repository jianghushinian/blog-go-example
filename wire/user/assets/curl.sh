curl --location --request POST 'http://127.0.0.1:8000/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "jianghushinian007@outlook.com",
    "nickname": "江湖十年",
    "username": "jianghushinian",
    "password": "pass"
}'
