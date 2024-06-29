build_dockerfile(){
    echo "FROM $1 as baseline" > Dockerfile
    echo "CMD echo '127.0.0.1 local-msn' >> /etc/hosts; /bin/sh -c 'while :; do sleep 1; done'" >> Dockerfile
}


build_dockerfile alpine:latest 
podman build -t my-build:latest -f Dockerfile .