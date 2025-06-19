# To generate keys....

```sh
# openssl genpkey -algorithm RSA -out production-key.pem \
    -pkeyopt rsa_keygen_bits:16384 \
    -pkeyopt rsa_keygen_primes:5
# openssl rsa -in production-key.pem -pubout > production-pubkey.pem
```