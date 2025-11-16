# dynamic-pod-init

A small helper container to initialize dynamic pods.

## Usage
The container uses the following environment variables:

- `TARGET_DIR` : the target directory to extract files to. Defaults to `/`
- `TEMPLATE_DIR` : template directory whose contents will be copied to `TARGET_DIR`. Ignored if empty
- `PATCHES` : a space-separated list of archive-file urls. These archives will be extracted to `TARGET_DIR`. Ignored if empty

## Example
This is an example k8s Pod configuration which will copy all contents from `/template` to `/app` and extract the provided archives to `/app`.
Note the shared volume `temp-app-data` which is necessary for the init-container's changes in `/app` to be persistent for the main application container.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: example-app
spec:
  volumes:
  - name: template-storage
    persistentVolumeClaim:
      claimName: template-pv-claim
  - name: temp-app-data
    emptyDir: {}
  initContainers:
  - name: dynamic-pod-init
    image: ghcr.io/benpueschel/dynamic-pod-init:main
    volumeMounts:
    - name: template-storage
      mountPath: /template
    - name: temp-app-data
      mountPath: /app
    env:
      - name: TARGET_DIR
        value: "/app"
      - name: TEMPLATE_DIR
        value: "/template"
      - name: "PATCHES"
        value: "https://example.com/patch.tar https://foo.bar/another_archive.zip"
  containers:
  - name: app
    image: example-image
    volumeMounts:
    - name: template-storage
      mountPath: /template
    - name: temp-app-data
      mountPath: /app
```
