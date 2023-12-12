# IMG FOR TESTING
FROM postgres:latest

ENV POSTGRES_USER=kasimov

ENV POSTGRES_PASSWORD=PasSwordsMan

ENV POSTGRES_DB=trivialDB

COPY storage/postgres/initTrivial.sql /docker-entrypoint-initdb.d/

EXPOSE 54432