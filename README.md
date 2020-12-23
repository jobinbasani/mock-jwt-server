### Mock JWT server 
Get a token - http://localhost:8994/token/userid  
Replace `userid` with any value  
JWKS - http://localhost:8994/.well-known/jwks.json  
Get token for a custom payload -
```
curl -X POST \
http://localhost:8994/token \
-d '{
"custom": "value",
"sub": "uid1"
}'
```