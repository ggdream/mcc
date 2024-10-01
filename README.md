# Mini CICD

将GitHub、GitLab、Gitea等服务的Webhook打到该项目，如果是`push`消息，则会拉下对应仓库代码，并按照配置文件进行cicd。整个过程都在部署该服务的主机完成。
无容器依赖，足够轻量，方便低成本、少人力的开发流程。

- type: server

用来打包并部署不退出的服务，比如web server、Job等。

~~~yaml
type: server
steps:
  - go mod download
  - go build -o build/app
apply: ./build
command: ["./app", "-c", "./config.yaml"]
~~~

- type: static

用来打包并部署静态资源，比如前端打包后的产物

~~~yaml
type: static
steps:
  - pnpm install
  - pnpm run build
apply: ./dist
~~~
