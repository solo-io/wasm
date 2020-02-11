FROM l.gcr.io/google/bazel:1.2.1

RUN apt update && \
    apt install bzip2 libxml2 nodejs -y

COPY build-filter.sh /build-filter.sh

ENTRYPOINT /bazel-build-filter.sh
