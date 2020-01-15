FROM golang:1.13

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN apt-get update
RUN apt-get -y install netcat openssh-server acl
RUN useradd -m sftpuser
RUN mkdir /home/sftpuser/.ssh
RUN echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDYl+tz9UPuQf7lu+tt4+F2ndrUbu4X2FlI0WmdEaDYgGA+G76nRN19cJK/NIXWH3DJWJG6lmZq1GwBLRrRWvKvqeiLN2Gn0df+V+POoMSYaokoLwYU4OHUOfavd7kY7AvPmEV/as3SG+ApiEZy7paYDUID98+HTcJ02j7Jjw/lecDutUySZkrIhqGVeT2F+sPgtnJvAEYJRmU/UqqOwSYbL3nFn3zH+d+yqpKuQpe/3GwWFIgzS7+gaHH3MTmTY3DG5OWRJLTvewKUSYhBJDGxaBH6+g9DjG0ZQpirQ2Y5uxmKcBpcMNd6f5GJE1aIhYk1E0p/ZpQsdNmOCky11lPB laboratorio@laboratorio-NE" >> /home/sftpuser/.ssh/authorized_keys
RUN echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCr01QfxA4oFHp1aJXyg8v9qwuZA6eUNuzifJj84GER4DsiwJJRV8Y0klMHHAIGovmgTkV0lkQzs6QE3AJuRDk1EiI8F1cDy+kB5QHpmHE+LZphEQMp/+aIu1SK1dBwJvg5J96Ys1/It1B9YRZMdJ+suu3qa2T5yGdXhR2McIIhiG2DPebX7hRekxf3avd0pI8buMQxElNFArkOEh1DBkn1qqbm3ZKoc5jLCWwGPA6C/s7Asmhb41Mu1//jOHNehCi5ST5P4Wjb4Bu2p92obV6r1HwK4G/TLXKnFP7blfRzAYIrBPn3TW3D0wBczb+Ao7StmIOecqzeM9yr6j5vwpYt duma@localhost.localdomain" >> /home/sftpuser/.ssh/authorized_keys
RUN chown -R sftpuser.sftpuser /home/sftpuser/.ssh
RUN chmod 640 /home/sftpuser/.ssh/authorized_keys

VOLUME /data/all

COPY css /data/all/css

EXPOSE 8000/tcp

ENTRYPOINT ["./entrypoint.sh"]
CMD ["app"]
