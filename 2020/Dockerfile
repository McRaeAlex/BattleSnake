FROM rust as builder
WORKDIR /usr/src/battlesnake
COPY . .
RUN cargo install --path .
EXPOSE 8000
CMD ["battlesnake-2020"]
