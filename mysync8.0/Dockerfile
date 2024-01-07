FROM ubuntu:22.04

RUN apt-get update && \
	apt-get install mysql-server-8.0 -y

RUN mkdir -p /var/log/
RUN mkdir -p /var/run/mysync/

RUN touch /var/log/mysync.log
RUN touch /var/run/mysync/mysync.lock

COPY ./mysync /usr/bin/mysync
COPY ./start.sh /start.sh
COPY ./mysqld.cnf /etc/mysql/mysql.conf.d/mysqld.cnf

CMD [ "/bin/bash", "/start.sh" ]
