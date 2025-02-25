FROM alpine

# Create a new user 'str' with user ID 1001
RUN adduser -D -u 1001 str

# Copy files
# todo check the naming of the binary
COPY str-ap-internal /opt/str
COPY docs/ /opt/docs/

RUN touch /opt/.env

# Change ownership of the /opt directory to the new user
RUN chown -R str:str /opt

WORKDIR /opt

# Switch user
USER str

ENTRYPOINT ["/opt/str"]