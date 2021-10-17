FROM lambci/lambda:go1.x

USER root

RUN yum install git make gcc -y
RUN yum install curl -y
RUN curl -L https://omnitruck.chef.io/install.sh | bash -s -- -P chef-workstation
USER sbx_user1051
