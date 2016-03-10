FROM selenium/standalone-chrome
EXPOSE 10000
ADD ./standalone /standalone
CMD [ "/bin/bash", "-c", "/opt/bin/entry_point.sh & /standalone" ]