FROM ubuntu:22.04

RUN apt-get update && \
	apt-get install default-jre -y && \
    apt-get install zookeeper -y

COPY ./mysql_proxy /mysql_proxy
COPY ./start.sh /start.sh
COPY ./main.go /main.go
RUN chmod +x /start.sh
RUN chmod +x /mysql_proxy

CMD [ "/bin/bash", "-c", "/start.sh" ]
