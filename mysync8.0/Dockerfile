FROM ubuntu:22.04

RUN apt-get update && \
	apt-get install mysql-server-8.0 -y


RUN mkdir -p /var/log/
RUN touch /var/log/mysync.log

RUN mkdir -p /var/run/mysync
RUN touch /var/run/mysync/mysync.lock

COPY ./mysync /usr/bin/mysync
COPY ./mysync /var/run/mysync
COPY ./mysync-init /etc/init.d/mysync-init
COPY ./start.sh /start.sh
COPY ./mysqld.cnf /etc/mysql/mysql.conf.d/mysqld.cnf

RUN chmod +x /etc/init.d/mysync-init

CMD [ "/bin/bash", "/start.sh" ]
