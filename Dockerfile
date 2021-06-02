FROM alpine
RUN apk --no-cache add curl
ADD  orders /orders
ENTRYPOINT [ "/orders" ]