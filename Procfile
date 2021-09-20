hraftd0: ./bin/hraftd -bootstrap-id hraftd0 -data-dir /tmp/raft -haddr :11000 -id hraftd0 -raddr 127.0.0.1:12000
hraftd1: sleep 5 && ./bin/hraftd -bootstrap-id hraftd0 -data-dir /tmp/raft -haddr :12000 -id hraftd1 -join :11000 -raddr 127.0.0.1:12001
hraftd2: sleep 15 && ./bin/hraftd -bootstrap-id hraftd0 -data-dir /tmp/raft -haddr :13000 -id hraftd2 -join :11000 -raddr 127.0.0.1:12002