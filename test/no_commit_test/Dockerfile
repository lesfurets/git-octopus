FROM lesfurets/octopus-tests:latest
ADD test.sh /home/
ADD bin /home/
RUN chmod +x /home/test.sh
WORKDIR /home/octopus-tests/
CMD /home/test.sh
