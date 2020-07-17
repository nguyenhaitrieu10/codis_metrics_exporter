# Codis metrics exporter
Export codis metrics to Prometheus
Convert metrics from codis proxy to Prometheus format

start by docker:

    $ docker-compose up -d

or start as service:

    $ cp codis-metrics-exporter.service /etc/systemd/system/
    
    $ mkdir -p /opt/codis-metrics-exporter
    
    $ cp main /opt/codis-metrics-exporter/main
    
    $ service codis-metrics-exporter start

![alt text](https://github.com/nguyenhaitrieu10/codis_metrics_exporter/blob/master/picture.png?raw=true)




