### Setup

1. in ```.env``` you can pick whether you want to use in-memory storage or PostgresSQL
2. ```Makefile``` provides commands for running app with both options

### Running the app

1. ```git clone https://github.com/DimaGitHahahab/ozon-fintech-posts.git```
2. navigate to the project root

#### With Postgres container

3. ```make compose-up``` to build and run
4. ```make compose-down``` to clean up

#### Without Postgres container

3. ```make docker-build``` to build image
4. ```make docker-run``` to run
5. ```make clean``` to clean up

### Connection

You can connect to GraphQL to see the schema and run queries at ```localhost:8080/root```

### Note

In Postgres option, by default there are some mock posts and comments being added in
migrations [here](https://github.com/DimaGitHahahab/ozon-fintech-posts/tree/main/migrations). 
