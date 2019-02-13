# docker build -t fareoffice/rig:latest .
# docker push fareoffice/rig:latest
FROM node:10.15.0-alpine

RUN mkdir -p /opt
ENV NODE_PATH /opt/app

WORKDIR /opt/app
COPY ./package.json .
COPY ./yarn.lock .

COPY . /opt/app
RUN yarn --production

ENTRYPOINT ["bin/rig"]
