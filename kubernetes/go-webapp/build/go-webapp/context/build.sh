cp ../../../src/kube-web build/kube-web
echo 'FROM registry.redhat.io/ubi9:latest as BASE' > Containerfile 
echo 'COPY build/kube-web /usr/bin/kube-web &&\ 
    chmod z+x /usr/bin/kube-web' >> Containerfile 
echo 'ENTRYPOINT /usr/bin/kube-web' >> Containerfile
podman build -f Containerfile  -t localhost/kube-webapp:latest .
