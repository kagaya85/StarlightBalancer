cpu rate
```
1-sum(increase(node_cpu_seconds_total{instance="10.244.1.2:9100",mode="idle"}[1m]))/sum(increase(node_cpu_seconds_total{instance="10.244.1.2:9100"}[1m]))

1-(avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[1m]))) 
```
