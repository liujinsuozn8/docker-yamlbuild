FROM centos

# 添加 kafak 安装包
ARG KAFKA_VERSION
ADD src/add/kafka_$KAFKA_VERSION /opt/module/
ENV KAFKA_HOME /opt/module/kafka
ENV PATH $PATH:$KAFKA_HOME/bin

# 引入 kafka 的 main.sh
ADD src/entrypoint/main.sh /kafka/base/main.sh

RUN chmod a+x /kafka/base/main.sh

VOLUME [ "/kafkaConfig" ]