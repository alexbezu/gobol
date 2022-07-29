docker run --rm -d --network mynet -p 3322:3322 --hostname db alexbezuglyi/mqimmudb

docker run --rm -d --network mynet -p 23567:23567 -v /home/oleksii/plexer/transactions:/app/transactions gobol