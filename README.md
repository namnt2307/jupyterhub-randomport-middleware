# Jupyterhub random port middleware

```
This api is use to get free port for Jupyterhub on Kubernetes when single user use hostNetwork: "true"
```
### How to use
````
curl -X POST "http://localhost:2307/getSpawnNode" -H "Content-type: application/json" -d '{"namespace": "default","podName":"namnt96","nodeSelector":"cpu","cpuLimit":"500m","cpuRequest":"200m","memoryRequest":"200M","memoryLimit":"500M"}'
````
### Sample response
````
{"podName":"namnt96","nodeSelector":"cpu","hostIP":"172.27.11.166","hostName":"node2","hostPort":0}
````
### Jupyterhub configuration
````
random_port: |
      import requests
      import json
      import socket
      async def custom_pre_spawn_hook(spawner):
        nodeSelector = {
          "gpu-server": "gpu",
          "default-dev-ml": "cpu"
          }
        username = spawner.user.name
        print(spawner.user_options)
        option=nodeSelector[spawner.user_options['profile']]
        print("Pod will be scheduled on",option,"server")
        header={'Content-Type':'Application/json'}
        data={
          'namespace': 'jupyterhub',
          'podName': 'namnt96',
          'nodeSelector': option,
          'cpuLimit':'500m',
          'cpuRequest':'200m',
          'memoryRequest':'200M',
          'memoryLimit':'500M'
        }

        host = requests.post(
          headers=header,
          url='http://dev.bigdata.local/dev-jupyterhub-freeport/v2/getSpawnNode',
          json=data,
        )
        hostName=host.json()['hostName']
        hostIP=host.json()['hostIP']
        hostPort=host.json()['hostPort']
        print(f"Pod will be spawned on {hostName}, ip: {hostIP} and port: {hostPort}")

        spawner.port=int(hostPort)

        spawner.node_selector={'kubernetes.io/hostname': hostName}
        print("node selector", spawner.port,spawner.node_selector)
      c.KubeSpawner.pre_spawn_hook = custom_pre_spawn_hook
````
