## Getting started

### Docker and Docker Compose

This project was created using the following versions:
Docker version 20.10.21
Docker Compose version v2.12.2

### Install node modules

If you want to be able to recompile and launch the web server on file change then this project requires Node and NPM.

Install the required modules by running the following command in the root directory:

```sh
npm i
```

### Tasks

There are a number of tasks specified in the .vscode/tasks.json file that can automate various things like starting the docker containers and the web server at the same time.

To run a task or collection of tasks open the command pallet (On a Mac the shortcut is `cmd + shift + P`) type in "run task" and then select the task you want to run.

Pay special attention to the following tasks as they are useful for getting the project up and running:

- "Start development"
- "Init outbox DB"

### Initialize the DB

Start the database docker container by either running `docker-compose up -d` in your terminal (omit the `-d` if you don't want detached mode). Or by running either the "Run docker containers" task or "Start development" task (see the tasks section for more information).

When the DB container is running run the "Init outbox DB" task.

You can use the following query to see some of the seed data:

```sql
select
	o.id,
  o.status,
  o.client_id,
  i.name,
  i.description,
  i.price
from
	orders o
join order_items oi on o.id = oi.order_id
join items i on i.id = oi.item_id
where
  o.client_id = 2;
```

### Start development

To get started with development run the task called "Start development". This will run the docker containers in the docker-compose.yml and start the web server.

If for some reason you just want to start the web server without the docker containers you can run the following command in the root directory:

```sh
npm run dev
```

### Connect to the DB from the terminal

If you want to connect to the DB to run some arbitrary queries or commands using the terminal you can get started using this:

```sh
docker exec -it outboxdb psql -d outbox -U postgres
```

### Removing the project containers/volumes

You can remove all the project containers but running this command in the root directory:

```sh
docker-compose rm -v
```

To remove the outbox DB volume - for example if you want to re-init the db - you can run the following command:

```sh
docker volume rm outbox_outboxdb_data
```
