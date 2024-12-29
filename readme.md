# gofiber-postgres-go

A book store API using Go, Fiber, and Postgres.

## Postgres Setup

1. **Install PostgreSQL using Homebrew:**

    ```sh
    brew install postgresql
    ```

2. **Verify the installation by checking the PostgreSQL version:**

    ```sh
    psql --version
    ```

3. **Start the PostgreSQL service:**

    ```sh
    brew services start postgresql
    ```

4. **Access the PostgreSQL command line interface (CLI) as the default `postgres` user:**

    ```sh
    psql -U postgres -d postgres
    ```

5. **Create a new user `dbuser` with a password:**

    ```sql
    CREATE USER dbuser WITH PASSWORD 'password';
    ```

6. **Create a new database named `books`:**

    ```sql
    CREATE DATABASE books;
    ```

7. **Grant all privileges on the `books` database to the user `dbuser`:**

    ```sql
    GRANT ALL PRIVILEGES ON DATABASE books TO dbuser;
    ```
