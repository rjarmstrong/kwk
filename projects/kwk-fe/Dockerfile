FROM nginx:1.13.3-alpine

COPY ./public/ /usr/share/nginx/html/

# MINIFY
RUN mkdir /build
COPY ./scripts/ /build
RUN /build/minify.sh
RUN rm /usr/share/nginx/html/index.html
COPY ./nginx.conf /etc/nginx/nginx.conf
