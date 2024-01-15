stream {
    upstream galera {
        server wos-0.wos-mysql-cluster.mysql.svc.cluster.local:3306;
    }

    server {
        listen 3306;
        proxy_pass galera;
        proxy_connect_timeout 1s;
    }
}

./zkCli.sh -server 10.244.120.76:2181 get /mysql/master | tail -n 2 | head -n 1 | tr -d '"'
