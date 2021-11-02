FROM centos:7
COPY hello /root/server
EXPOSE 8080
CMD /root/server
