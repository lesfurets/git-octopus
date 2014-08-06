FROM lesfurets/octopus-tests:merge_failure_latest
ADD test.sh /home/
ADD bin /home/
RUN chmod +x /home/test.sh
WORKDIR /home/octopus-tests/
CMD /home/test.sh