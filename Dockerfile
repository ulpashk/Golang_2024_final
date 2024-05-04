FROM golang:latest
RUN mkdir /build
WORKDIR /build
RUN cd /build && rm -rf goproject
RUN cd /build && git clone https://github.com/ulpashk/Golang_2024_final.git
RUN cd /build/goproject && go build
EXPOSE 8080
ENTRYPOINT [ "/build/goproject/goproject" ]

