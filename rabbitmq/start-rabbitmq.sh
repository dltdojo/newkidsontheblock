docker run -d -p 5672:5672 -p 15672:15672 -v /srv/raw/virtualisation/rabbitmq-docker-storage/log:/data/log -v /srv/raw/virtualisation/rabbitmq-docker-storage/mnesia:/data/mnesia --name=rabbitmq_blockchain ralph/rabbitmq:latest