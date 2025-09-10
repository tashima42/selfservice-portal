# selfservice-portal

Use this with extreme caution. To securely run this, you must follow these recommendations:

1. Never run this publicly or without an external authentication method.
2. Make sure that your reverse proxy overwrites the `X-Real-IP` header as this might lead to CSRF attacks or be the entry point to a more sofisticated Pangolin attack.

Future improvements:

1. Proper error handling
2. Duplication checks
3. Pagination
4. Filters
