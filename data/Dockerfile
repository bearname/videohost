FROM migrate/migrate:v4.11.0
COPY ./sql /mysql/migration

# override base image ENTRYPOINT
ENTRYPOINT ["/usr/bin/env"]
# "sh -c" is required to extract ENV values using shell
CMD ["sh", "-c", "./migrate --path=/migrations --database \"mysql://$DB_USER:$DB_PASSWORD@($DB_HOST:$DB_PORT)/$DB_NAME\" up"]