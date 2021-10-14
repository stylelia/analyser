FROM lambci/lambda:go1.x

USER root

RUN yum install git -y

USER sbx_user1051