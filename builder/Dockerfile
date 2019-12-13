FROM l.gcr.io/google/bazel:1.2.1

RUN apt update && \
    apt install bzip2 libxml2 -y

COPY bazel-build-filter.sh /bazel-build-filter.sh

ENTRYPOINT /bazel-build-filter.sh
