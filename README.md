
````
curl -X POST "http://localhost:8000/api/getSpawnNode" -H "Content-type: application/json" -d '{"namespace": "default","podName":"namnt96","nodeSelector":"cpu","cpuLimit":"500m","cpuRequest":"200m","memoryRequest":"200M","memoryLimit":"500M"}'
````
