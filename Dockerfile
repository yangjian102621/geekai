# FROM 表示设置要制作的镜像基于哪个镜像，FROM指令必须是整个Dockerfile的第一个指令，如果指定的镜像不存在默认会自动从Docker Hub上下载。
FROM centos:7 

WORKDIR /usr/src/app
COPY src/bin/wechatGPT-amd64-linux /usr/src/app

# 容器对外暴露的端口号，这里和配置文件保持一致就可以
EXPOSE 5678

# 容器启动时执行的命令
CMD ["./wechatGPT-amd64-linux"]
