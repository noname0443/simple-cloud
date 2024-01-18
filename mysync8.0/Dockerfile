FROM ubuntu:22.04

RUN apt-get update && \
	apt-get install mysql-server-8.0 -y

RUN mkdir -p /var/log/
RUN mkdir -p /var/run/mysync/

RUN touch /var/log/mysync.log
RUN touch /var/run/mysync/mysync.lock

RUN echo 10 > /tmp/usedspace
RUN echo false > /tmp/readonly
RUN rm /etc/mysql/mysql.conf.d/mysqld.cnf

COPY ./mysync_fork /usr/bin/mysync
RUN chmod +x /usr/bin/mysync
