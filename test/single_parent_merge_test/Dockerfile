FROM lesfurets/octopus-tests:simple_merge_latest
ADD test.sh /home/
ADD bin /usr/local/bin
RUN chmod +x /home/test.sh
#RUN git branch -D features/feat2 features/feat3
#RUN git branch qa
WORKDIR /home/octopus-tests/
CMD /home/test.sh
