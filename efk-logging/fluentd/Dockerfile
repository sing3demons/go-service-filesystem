# Use the Fluentd base image
FROM fluent/fluentd:v1.12.0-debian-1.0

# Switch to root user for installation
USER root
# Update system packages and install development tools
RUN apt-get update && apt-get install -y build-essential && apt-get clean
# Install an older version of faraday compatible with Ruby 2.6
RUN gem install faraday -v 2.8.1
RUN gem install elasticsearch -v 7.13.0
RUN gem install elasticsearch-transport -v 7.13.0

RUN gem install fluent-plugin-elasticsearch --no-document -v '5.0.3'
# Switch back to the fluent user
USER fluent
