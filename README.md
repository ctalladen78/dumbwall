Yet another failure

# dumbwall
HackerNews like toy app to post whatever you want and rate other posts

## pre install

_Export database source name to public env variable:_
```
export DSN=postgres://postgres:mysecretpassword@localhost:5432/dumbwall?sslmode=disable
```

_Run database migrations:_
```
make prepare
make migrate
```

_Generate priv/pub keys for JWT signing:_
```
ssh-keygen -t rsa -b 4096 -f private.key
openssl rsa -in private.key -pubout -outform PEM -out public.key
```
