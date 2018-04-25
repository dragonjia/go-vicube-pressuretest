# go-vicube-pressuretest


##背景

	利用golang 充分利用打压源性能，对restful 接口进行打压力，同时收集性能数据。
	
##功能清单

- 可定制打压并发进程数
- 可收集性能数据 （状态、各个环节响应时长）

##tudo

- 能够输出性能报告
- 可配置打压时长和分节奏策略、以及每个节奏的压力定义



##性能数据情况
######当前:2018-4-25
~~~
200 OK 
DNS lookup:           0 ms
TCP connection:       1 ms
TLS handshake:        0 ms
Server processing:    9 ms
Content transfer:     0 ms

Name Lookup:       0 ms
Connect:           1 ms
Pre Transfer:      1 ms
Start Transfer:   10 ms
Total:            10 ms
~~~








