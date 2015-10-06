FROM debian
ADD ./znode_size /work/
WORKDIR /work/
ENV ZK=master.mesos:2181
ENV DIR=/universe/mon-marathon-service/state
ENV MIN=1024
CMD ./znode_size -zk=$ZK -dir=$DIR -min-sz=$MIN; sleep 10
