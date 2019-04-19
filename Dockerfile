FROM alpine:3.9
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    iptables
ADD build/bin/linux/kube2ram /bin/kube2ram
ENTRYPOINT ["kube2ram"]
