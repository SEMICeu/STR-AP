# ingress-nginx

This chart has been made based on the raw yaml file:
https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.1/deploy/static/provider/aws/deploy.yaml

important changes:
* Adding `- --enable-ssl-passthrough` in the deployment container args.
* The ingress should have those annotations
```
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/backend-protocol: "HTTPS"
    nginx.ingress.kubernetes.io/ssl-passthrough: "true"
```
* There's only 1 listener configured on TCP port 443 with passing a certificate; we only want to do SSL termination in the application itself!
