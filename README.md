# streaming-example
 An example web server that streams audio and video with golang


## Docker
Build the docker image: `docker build --tag {image_name} .`

Confirm that the image is built: `docker image ls`

Run the docker container locally: `docker run -p{localhost_port}:8080 {image_name}` (optional `-d` tag to run detached)

Stop docker container: `docker stop {container_id}`