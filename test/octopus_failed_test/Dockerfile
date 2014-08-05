FROM lesfurets/octopus-test_simple_fail
ADD test.sh /home/
ADD bin /home/
RUN chmod +x /home/test.sh
WORKDIR /home/octopus-tests/
CMD /home/test.sh