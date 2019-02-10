# docker build -t fareoffice/infra-cli:latest .
# docker push fareoffice/infra-cli:latest
FROM node:10.15.0-alpine

RUN mkdir -p /opt
ENV NODE_PATH /opt/app

WORKDIR /opt/app
COPY ./package.json .
COPY ./yarn.lock .

COPY . /opt/app
RUN yarn --production

ENTRYPOINT ["bin/cli"]
