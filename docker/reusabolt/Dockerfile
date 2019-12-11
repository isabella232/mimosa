FROM gcr.io/reusabolt/reusabolt-preview

RUN rm -rf /root/reusabolt/Boltdir/*
COPY Boltdir /root/reusabolt/Boltdir
RUN bolt puppetfile install
RUN bolt task show
