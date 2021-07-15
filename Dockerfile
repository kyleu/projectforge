FROM scratch
ENTRYPOINT ["/projectforge", "-a", "0.0.0.0"]
EXPOSE 14000
COPY projectforge /
