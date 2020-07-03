# GOLANG
docker container prune   使用後要回復ＤＢ資料

docker network create elastic_stack
# ES
<p>docker run -d --name elasticsearch --net elastic_stack -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.7.0<p>
<p>docker container start elasticsearch<p>

# kibana
<p>docker run -d --name kibana --net elastic_stack -p 5601:5601 docker.elastic.co/kibana/kibana:7.7.0<p>
<p>docker container start kibana<p>



# redis
<p>docker run -d --name redis-lab -p 6379:6379 redis<p>
<p>docker container start redis-lab<p>

# mysql
<p>docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:5.7.24<p>
<p>docker container start mysql<p>

# phpmyadmin
<p>docker run -d --name phpmyadmin --link mysql -e PMA_HOST="mysql" -p 8080:80 phpmyadmin/phpmyadmin<p>
<p>docker container start phpmyadmin<p>

# 建置

# 說明