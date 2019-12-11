#### Offers store
## Приложение для загрузки офферов yml в cassandra

#### Приложение работает с двумя url:
1. **`http://localhost/loadShareFile?source=\\DESKTOP\Users\user\tmp\data\someFile.xml`** - загрузка файла из сетевой папки
2. **`http://localhost/loadUrl?source=https://some-host.ru/some-url/someFile.xml`** - загрузка web файла

Перед запуском приложения необходима cassandra.  Предлагается зпустить скрипт: ./queries/initDB. В данном файле создается keyspace и таблица offers.

####Настройки окружения в файле .env:

MAX_WORKERS=5 - количество воркеров учатсвующих при загрузке в cassandra. Влияет на стабильность вместе с BATCH_SIZE

BATCH_SIZE=250 - размер отправляемой пачки запроса в cassandra. Влияет на стабильность вместе с MAX_WORKERS

DB_HOST=localhost - расположение cassandra

DB_PORT=9040 - порт cassandra

DB_USER=cassandra - login cassandra

DB_PASSWORD=cassandra - password cassandra

Keyspace=offers - Keyspace cassandra, возможен иной если не использовали в предлагаемом скрипте

Consistency=One - режим Consistency cassandra, по умолчанию One

CLUSTER_TIMEOUT=120 - timeout после чего приложение считает что запрос не успешен и пишет об этом в лог, но как показала практика - если отдали пакет, то он будет загружен

####Сборка проекта:

docker build -t offersstore .

docker-compose up --build

#### Для теста:

## Поднимаем cassandra локально:

docker network create testnetwork

docker run -p 9042:9042 --name cassandra --network testnetwork -d cassandra:3.11.5

docker run --name cassandra2 --network testnetwork -d -e CASSANDRA_SEEDS=cassandra cassandra:3.11.5

Проверяем работу cassandra

docker run -it --network testnetwork --rm cassandra cqlsh cassandra
