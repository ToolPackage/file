# file

- 启动server：`docker-compose up`
- 在启动之前构建image：`docker-compose up --build`
- 在启动之前重新创建容器：`docker-compose up --force-recreate`
- 守护模式启动server：`docker-compose up -d`
- 重启server：`docker-compose restart`
- 停止server：`docker-compose stop`
- 杀掉server：`docker-compose kill`
- 查看server状态：`docker-compose ps`
- 删掉server容器：`docker-compose rm`
- 清理多阶段构建的中间镜像：`docker image prune --filter label=stage=builder`
