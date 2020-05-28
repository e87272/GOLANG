# GOLANG

# ES
<p>docker run -d -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.7.0<p>

# redis
<p>docker run -d --name redis-lab -p 6379:6379 redis<p>

# mysql
<p>docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:5.7.24<p>

# phpmyadmin
<p>docker run -d --name phpmyadmin --link mysql -e PMA_HOST="mysql" -p 8080:80 phpmyadmin/phpmyadmin<p>

# 建置

# 說明