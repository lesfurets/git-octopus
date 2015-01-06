FROM lesfurets/octopus-tests:simple_merge_latest
ADD test.sh /home/
ADD bin /usr/local/bin
RUN chmod +x /home/test.sh
WORKDIR /home/octopus-tests/
CMD /home/test.sh
