############################################################
# Build stage
############################################################






############################################################
# Run stage
############################################################

FROM alpine:latest

# install binary and remove cache
RUN apk update \
#    && apk upgrade \
    && apk add iperf3
#    && rm -rf /var/lib/apt/lists/*

# Expose the default iperf3 server port
#EXPOSE 5201

# entrypoint allows you to pass your arguments to the container at runtime
# very similar to a binary you would run. For example, in the following
# docker run -it <IMAGE> --help' is like running 'iperf3 --help'
#ENTRYPOINT ["iperf3"]
