FROM postgres:17
RUN apt-get update && apt-get -y install postgresql-17-cron
COPY ./initdb.d/ /docker-entrypoint-initdb.d/
