# docker-yamlbuild
- 基于yaml文件批量编译docker镜像

## 编译指令
- 编译，并生成可执行文件 `docker-yamlbuild`
    ```sh
    go build
    ```

## 镜像编译示例

```sh
# 读取 ./test/test.yml 文件内的编译参数
# 从 test/img/ 中获取 Dockerfile 等编译所需的内容
# test/resource 保存编译时所需要的安装包，并在该路径下执行编译（编译前，会从 test/img/ 下进行拷贝） 

./docker-yamlbuild -y ./test/test.yml --img-dir test/img/ --build-dir test/resource
```

## 指令参数

|参数|功能|必须指定|
|-|-|-|
|-y|指定 yaml 参数|Yes|
|--img-dir| Dockerfile 所在的父路径 |Yes|
|--build-dir| 编译路径 |Yes|
|-o| 不执行编译，只输出编译时的指令 |No|


## 如何保存镜像编译所需的内容
- 编译时所需要的目录结构
    ```
    - ${--img-dir}
        - imageID
            - src
                - xxx.sh
                - yyy.txt
            - Dockerfile
    - ${--build-dir}
        - 编译时所需要的安装包等内容
    ```

- 每一个 image 除了 `Dockerfile`，额外的代码，如shell文件等，需要保存到 `src` 目录
- 如果有需要统一管理的资源，如 image 所需的安装包，可以保存到 `${--build-dir}` 所指定的路径下
- 如果要使用 `Dockerfile` 之外的文件做编译，需要在 yaml 文件中通过 `file` 来设置
- 如果某个 image 需要在特殊的目录中进行编译，需要在 yaml 文件中通过 `build-dir` 来设置

## yaml的编写
- yaml 示例
    ```yaml
    - id: test-img # image 在 ${--img-dir} 下的保存路径 
      build-arg:
        KAFKA_VERSION: 123
        bbb: 234
      tag:
        - T1
        - T2
    - id: test-img2
      build-arg:
        KAFKA_VERSION: 456
        bbb: 678
    ```

- 编写方式
    - 整体为一个list
    - list中的一个元素，为一个镜像编译时所需要的参数
    - 每一个镜像**必须指定 id**

- `docker build` 的指令参数所对应的yaml内容
    
    |docker参数|yaml参数|yaml类型|必须设置|
    |-|-|-|-|
    |--build-arg|build-arg|map|No|
    |-t|tag|list|No，如果没有设置，会替换为 `id`|

- 额外的yaml参数

    |yaml参数|yaml类型|必须设定|功能|
    |-|-|-|-|
    |id|key: value|Yes|标识image在 `--img-dir` 下的保存路径|
    |file|key: value|No|用于替换：`${--img-dir}/id/Dockerfile`，并且会在`file`所在的目录下，搜索 `ADD`、`COPY` 所使用的本地资源|
    |build-dir|key: value|No|用于替换指令参数 `${--build-dir}`|

## 镜像的编译流程

1. 执行 `docker-yamlbuild` 指令
2. 按照 yaml 文件中的 list 的顺序，依次解析每一个镜像的`参数`和`Dockerfile`
    - 解析 `Dockerfile` 时，会抽取 `ENV`、`ARG`、`COPY`、`ADD`
    - `COPY`、`ADD` 中本地资源路径中的参数会被替换为 `ENV`、`ARG` 中的实际值，来得到真实的本地资源路径
    - **如果 `COPY`、`ADD` 指定的路径不存在，会抛出异常，并暂停编译**
        - 会同时在 `--build-dir`、`${--img-dir}/id/` 两个路径下搜索
3. 将 `Dockerfile` 、`src/` 拷贝到 `--build-dir` 路径下
4. 根据 `COPY`、`ADD` 创建 `.dockerignore` 文件，防止 `--build-dir` 路径下资源过多导致的编译速度下降
5. 当 `--build-dir` 路径下的资源准备完成之后，执行编译
6. 编译结束后，会将拷贝的资源删掉，包括：`Dockerfile` 、`src/`、`.dockerignore`

## 注意事项
- 无法解析 `ENV` 中包含多个 `KEY=VALUE` 形式的Step，可能会造成解析失败