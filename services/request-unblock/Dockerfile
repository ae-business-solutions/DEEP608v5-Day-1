FROM ekidd/rust-musl-builder:nightly-2021-02-13 AS rust-build
RUN rustup target add x86_64-unknown-linux-musl
# https://github.com/emk/rust-musl-builder#making-static-releases-with-travis-ci-and-github
RUN sudo chown -R rust:rust /home
RUN mkdir /home/source
WORKDIR /home/source
COPY ./Cargo.toml .
COPY ./Rocket.toml .
COPY ./src/ ./src/
COPY ./templates/ ./templates/
RUN cargo build --bins --release --target x86_64-unknown-linux-musl
RUN strip ./target/x86_64-unknown-linux-musl/release/request-unblock

FROM scratch
COPY --from=0 /home/source/target/x86_64-unknown-linux-musl/release/request-unblock /
COPY --from=0 /home/source/Rocket.toml /
COPY --from=0 /home/source/templates/ /templates/
CMD ["./request-unblock"]