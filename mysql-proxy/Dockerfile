FROM ubuntu:22.04

RUN apt-get update && \
	apt-get install default-jre -y && \
    apt-get install zookeeper -y

COPY ./mysql_proxy /mysql_proxy
COPY ./sample.sh /sample.sh
RUN chmod +x /sample.sh

CMD [ "/bin/bash", "-c", "/mysql_proxy" ]