# go test

docker run -p 13306:3306 --name mysql -v ~/docker/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=Broadfun@123  -d mysql:5.7

systemctl daemon-reload

systemctl restart docker.service

grant all privileges on *.* to root@'%' identified by '123456';